export interface XMLName {
  Space: string;
  Local: string;
}

export interface Finished {
  time: number;
  time_str: string;
  elapsed: number;
  summary: string;
  exit: string;
  error_msg: string;
}

export interface HostStats {
  up: number;
  down: number;
  total: number;
}

export interface PortUsed {
  state: string;
  proto: string;
  port_id: number;
}

export interface OsMatch {
  name: string;
  accuracy: number;
  line: number;
  os_classes: OsClass[];
}

export interface OsClass {
  vendor: string;
  os_generation: string;
  type: string;
  accuracy: number;
  os_family: string;
  cpes: string[];
}

export interface OsFingerprint {
  fingerprint: string;
}

export interface HostAddress {
  addr: string;
  addr_type: string;
  vendor: string;
}

export interface Reason {
  reason: string;
  count: number;
}

export interface ExtraPort {
  state: string;
  count: number;
  reasons: Reason[];
}

export interface HostName {
  name: string;
  type: string;
}

export interface ServiceInfo {
  device_type: string;
  extra_info: string;
  high_version: string;
  hostname: string;
  low_version: string;
  method: string;
  name: string;
  os_type: string;
  product: string;
  proto: string;
  rpc_num: string;
  service_fp: string;
  tunnel: string;
  version: string;
  confidence: number;
  cpes: string[];
}

export interface PortState {
  state: string;
  reason: string;
  reason_ip: string;
  reason_ttl: number;
}

export interface Task {
  time: number;
  task: string;
  extra_info: string;
}

export interface TaskProgress {
  percent: number;
  remaining: number;
  task: string;
  etc: number;
  time: number;
}

export interface TaskEnd {
  time: number;
  task: string;
  extra_info: string;
}

export interface Host {
  distance: {
    value: number;
  };
  end_time: number;
  ip_id_sequence: {
    class: string;
    values: string;
  };
  os: {
    ports_used: PortUsed[];
    os_matches: OsMatch[];
    os_fingerprints: OsFingerprint[];
  };
  start_time: number;
  timed_out: boolean;
  status: {
    state: string;
    reason: string;
    reason_ttl: number;
  };
  tcp_sequence: {
    index: number;
    difficulty: string;
    values: string;
  };
  tcp_ts_sequence: {
    class: string;
    values: string;
  };
  times: {
    srtt: string;
    rttv: string;
    to: string;
  };
  trace: {
    proto: string;
    port: number;
    hops: null | unknown; // Adjust the type according to data pattern if needed
  };
  uptime: {
    seconds: number;
    last_boot: string;
  };
  comment: string;
  addresses: HostAddress[];
  extra_ports: ExtraPort[];
  hostnames: HostName[];
  host_scripts: null;
  ports: {
    id: number;
    protocol: string;
    owner: {
      name: string;
    };
    service: ServiceInfo;
    state: PortState;
    scripts: null;
  }[];
  smurfs: null;
}

export interface NmapRun {
  XMLName: XMLName;
  args: string;
  profile_name: string;
  scanner: string;
  start_str: string;
  version: string;
  xml_output_version: string;
  debugging: {
    level: number;
  };
  run_stats: {
    finished: Finished;
    hosts: HostStats;
  };
  scan_info: {
    num_services: number;
    protocol: string;
    scan_flags: string;
    services: string;
    type: string;
  };
  start: number;
  verbose: {
    level: number;
  };
  hosts: Host[];
  post_scripts: null;
  pre_scripts: null;
  targets: null;
  task_begin: Task[];
  task_progress: TaskProgress[];
  task_end: TaskEnd[];
  NmapErrors: null;
}
