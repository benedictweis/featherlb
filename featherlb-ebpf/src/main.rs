#![no_std]
#![no_main]

use core::mem;

use aya_ebpf::{
    bindings::xdp_action::{self, XDP_TX},
    macros::{map, xdp},
    maps::HashMap,
    programs::XdpContext,
};
use aya_log_ebpf::{debug, info};
use network_types::{
    eth::{EthHdr, EtherType},
    ip::{IpProto, Ipv4Hdr},
    tcp::TcpHdr,
    udp::UdpHdr,
};

static load_balancer_ip: u32 = 0x0AB48403;

#[map]
static REWRITES: HashMap<u16, u32> = HashMap::<u16, u32>::with_max_entries(1024, 0);

#[xdp]
pub fn featherlb(ctx: XdpContext) -> u32 {
    match try_featherlb(ctx) {
        Ok(ret) => ret,
        Err(_) => xdp_action::XDP_ABORTED,
    }
}

#[inline(always)]
fn ptr_at_impl<T>(ctx: &XdpContext, offset: usize) -> Result<usize, ()> {
    let start = ctx.data();
    let end = ctx.data_end();
    let len = mem::size_of::<T>();

    if start + offset + len > end {
        return Err(());
    }

    Ok(start + offset)
}

#[inline(always)]
fn ptr_at<T>(ctx: &XdpContext, offset: usize) -> Result<*const T, ()> {
    ptr_at_impl::<T>(ctx, offset).map(|p| p as *const T)
}

#[inline(always)]
fn ptr_at_mut<T>(ctx: &XdpContext, offset: usize) -> Result<*mut T, ()> {
    ptr_at_impl::<T>(ctx, offset).map(|p| p as *mut T)
}

fn rewrite_ip(_adress: u32, port: u16) -> Option<u32> {
    unsafe { REWRITES.get(&port).copied() }
}

fn try_featherlb(ctx: XdpContext) -> Result<u32, ()> {
    let ethhdr: *const EthHdr = ptr_at(&ctx, 0)?;
    match unsafe { (*ethhdr).ether_type } {
        EtherType::Ipv4 => {}
        _ => return Ok(xdp_action::XDP_PASS),
    }

    let ipv4hdr: *const Ipv4Hdr = ptr_at(&ctx, EthHdr::LEN)?;
    let source_addr = u32::from_be(unsafe { (*ipv4hdr).src_addr });
    let destination_addr = u32::from_be(unsafe { (*ipv4hdr).dst_addr });

    let (source_port, destination_port) = match unsafe { (*ipv4hdr).proto } {
        IpProto::Tcp => {
            let tcphdr: *const TcpHdr = ptr_at(&ctx, EthHdr::LEN + Ipv4Hdr::LEN)?;
            (
                u16::from_be(unsafe { (*tcphdr).source }),
                u16::from_be(unsafe { (*tcphdr).dest }),
            )
        }
        IpProto::Udp => {
            let udphdr: *const UdpHdr = ptr_at(&ctx, EthHdr::LEN + Ipv4Hdr::LEN)?;
            (
                u16::from_be(unsafe { (*udphdr).source }),
                u16::from_be(unsafe { (*udphdr).dest }),
            )
        }
        _ => return Err(()),
    };

    debug!(
        &ctx,
        "SRC IP: {:i}, SRC PORT: {}, DST IP: {:i}, DST PORT: {}",
        source_addr,
        source_port,
        destination_addr,
        destination_port
    );

    if let Some(new_ip) = rewrite_ip(destination_addr, destination_port) {
        info!(
            &ctx,
            "Rewriting DST IP to {:i} for SRC IP: {:i} and SRC PORT: {}",
            new_ip,
            destination_addr,
            destination_port
        );

        // Rewrite the destination IP address
        unsafe {
            let ipv4hdr_mut: *mut Ipv4Hdr = ptr_at_mut(&ctx, EthHdr::LEN)?;
            (*ipv4hdr_mut).dst_addr = u32::to_be(new_ip);

            // Update IP source address to the load balancer's IP
            (*ipv4hdr_mut).src_addr = u32::to_be(load_balancer_ip);

            // Recalculate IP checksum
            (*ipv4hdr_mut).check = 0;
            (*ipv4hdr_mut).check = compute_ip_checksum(ipv4hdr_mut);
        }

        // Update Ethernet source MAC address to the current lb's MAC
        let ethhdr_mut: *mut EthHdr = ptr_at_mut(&ctx, 0)?;
        unsafe {
            (*ethhdr_mut).h_source = load_balancer_mac;
            (*ethhdr_mut).h_dest = backend_mac;
        }

        return Ok(XDP_TX);
    } else {
        info!(&ctx, "Packet from backend, routing correctly");

        unsafe {
            let ipv4hdr_mut: *mut Ipv4Hdr = ptr_at_mut(&ctx, EthHdr::LEN)?;
            (*ipv4hdr_mut).dst_addr = u32::to_be(client_ip);

            // Update IP source address to the load balancer's IP
            (*ipv4hdr_mut).src_addr = u32::to_be(load_balancer_ip);

            // Recalculate IP checksum
            (*ipv4hdr_mut).check = 0;
            (*ipv4hdr_mut).check = compute_ip_checksum(ipv4hdr_mut);
        }

        // Update Ethernet source MAC address to the current lb's MAC
        let ethhdr_mut: *mut EthHdr = ptr_at_mut(&ctx, 0)?;
        unsafe {
            (*ethhdr_mut).src_addr = load_balancer_mac;
            (*ethhdr_mut).dst_addr = client_mac;
        }

        return Ok(XDP_TX);
    }

    Ok(xdp_action::XDP_PASS)
}

#[cfg(not(test))]
#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}
