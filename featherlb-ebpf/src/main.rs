#![no_std]
#![no_main]

use core::mem;

use aya_ebpf::{
    bindings::xdp_action::{self, XDP_TX},
    helpers::gen::bpf_csum_diff,
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

const LOAD_BALENCER_IP_ADDR: u32 = (10 << 24) | (180 << 16) | (132 << 8) | 2;
const LOAD_BALENCER_MAC_ADDR: [u8; 6] = [0x02, 0x42, 0x0a, 0xb4, 0x84, 0x02];
const CLIENT_IP_ADDR: u32 = (10 << 24) | (180 << 16) | (132 << 8) | 3;
const CLIENT_MAC_ADDR: [u8; 6] = [0x02, 0x42, 0x0a, 0xb4, 0x84, 0x03];
const BACKEND_IP_ADDR: u32 = (10 << 24) | (180 << 16) | (132 << 8) | 4;
const BACKEND_MAC_ADDR: [u8; 6] = [0x02, 0x42, 0x0a, 0xb4, 0x84, 0x04];

#[map]
static REWRITES: HashMap<u16, u32> = HashMap::<u16, u32>::with_max_entries(1024, 0);

#[cfg(not(test))]
#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}

#[xdp]
pub fn featherlb(ctx: XdpContext) -> u32 {
    match try_featherlb(ctx) {
        Ok(ret) => ret,
        Err(_) => xdp_action::XDP_ABORTED,
    }
}

fn try_featherlb(ctx: XdpContext) -> Result<u32, ()> {
    let ethhdr: *mut EthHdr = ptr_at_mut(&ctx, 0)?;
    match unsafe { (*ethhdr).ether_type } {
        EtherType::Ipv4 => {}
        _ => return Ok(xdp_action::XDP_PASS),
    }

    let source_mac_addr = unsafe { (*ethhdr).src_addr };
    let destination_mac_addr = unsafe { (*ethhdr).dst_addr };

    let ipv4hdr: *mut Ipv4Hdr = ptr_at_mut(&ctx, EthHdr::LEN)?;
    let source_ip_addr = u32::from_be(unsafe { (*ipv4hdr).src_addr });
    let destination_ip_addr = u32::from_be(unsafe { (*ipv4hdr).dst_addr });

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

    let source_mac_addr_print: &[u8] = &source_mac_addr;
    let destination_mac_addr_print: &[u8] = &destination_mac_addr;

    debug!(
        &ctx,
        "SRC: {:x}, {:i}, {} DST: {:x}, {:i}, {}",
        source_mac_addr_print,
        source_ip_addr,
        source_port,
        destination_mac_addr_print,
        destination_ip_addr,
        destination_port
    );

    if source_ip_addr == CLIENT_IP_ADDR {
        info!(&ctx, "Client sent a packet");

        unsafe {
            (*ethhdr).dst_addr[5] = BACKEND_MAC_ADDR[5];
            (*ipv4hdr).dst_addr = BACKEND_IP_ADDR.to_be();
        }
    } else if source_ip_addr == BACKEND_IP_ADDR {
        info!(&ctx, "Backend sent a packet");

        unsafe {
            (*ethhdr).dst_addr[5] = CLIENT_MAC_ADDR[5];
            (*ipv4hdr).dst_addr = CLIENT_IP_ADDR.to_be();
        }
    } else {
        info!(&ctx, "Unrelated packet");

        return Ok(xdp_action::XDP_PASS);
    }

    unsafe {
        (*ethhdr).src_addr[5] = LOAD_BALENCER_MAC_ADDR[5];
        (*ipv4hdr).src_addr = u32::from_be(LOAD_BALENCER_IP_ADDR);
    }

    let full_cksum = unsafe {
        bpf_csum_diff(
            mem::MaybeUninit::zeroed().assume_init(),
            0,
            ipv4hdr as *mut u32,
            Ipv4Hdr::LEN as u32,
            0,
        )
    } as u64;

    unsafe {
        info!(&ctx, "{}", (*ipv4hdr).check);
    }

    unsafe { (*ipv4hdr).check = csum_fold_helper(full_cksum) };

    unsafe {
        info!(&ctx, "{}", (*ipv4hdr).check);
    }

    Ok(xdp_action::XDP_TX)
}

#[inline(always)]
fn ptr_at<T>(ctx: &XdpContext, offset: usize) -> Result<*const T, ()> {
    ptr_at_impl::<T>(ctx, offset).map(|p| p as *const T)
}

#[inline(always)]
fn ptr_at_mut<T>(ctx: &XdpContext, offset: usize) -> Result<*mut T, ()> {
    ptr_at_impl::<T>(ctx, offset).map(|p| p as *mut T)
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
pub fn csum_fold_helper(mut csum: u64) -> u16 {
    for _i in 0..4 {
        if (csum >> 16) > 0 {
            csum = (csum & 0xffff) + (csum >> 16);
        }
    }
    return !(csum as u16);
}
