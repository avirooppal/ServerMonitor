export interface LoadAvg {
    load1: number;
    load5: number;
    load15: number;
}

export interface MemoryInfo {
    total: number;
    available: number;
    used: number;
    usedPercent: number;
    free: number;
    buffers: number;
    cached: number;
}

export interface SwapInfo {
    total: number;
    used: number;
    free: number;
    usedPercent: number;
}

export interface HostInfo {
    hostname: string;
    uptime: number;
    bootTime: number;
    procs: number;
    os: string;
    platform: string;
    platformFamily: string;
    platformVersion: string;
    kernelVersion: string;
    kernelArch: string;
    virtualizationSystem: string;
    virtualizationRole: string;
    hostId: string;
}

export interface ContainerInfo {
    id: string;
    name: string;
    image: string;
    state: string;
    status: string;
    created: number;
    cpu_percent: number;
    memory_usage: number;
    memory_limit: number;
}

export interface DiskInfo {
    path: string;
    total: number;
    used: number;
    free: number;
    used_percent: number;
    read_rate: number;
    write_rate: number;
}

export interface NetInterface {
    name: string;
    recv_rate: number;
    sent_rate: number;
}

export interface NetworkStats {
    interfaces: NetInterface[];
    total_recv: number;
    total_sent: number;
}

export interface ProcessInfo {
    pid: number;
    name: string;
    cpu: number;
    mem: number;
    username: string;
}

export interface SystemMetrics {
    cpu: number[];
    cpu_total: number;
    load_avg: LoadAvg;
    memory: MemoryInfo;
    swap: SwapInfo;
    disks: DiskInfo[];
    network: NetworkStats;
    processes: ProcessInfo[];
    containers: ContainerInfo[];
    host_info: HostInfo;
    last_update: string;
}
