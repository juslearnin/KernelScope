import { useEffect, useMemo, useState } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  type Node,
  type Edge,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import "./App.css";

type FileDescriptor = {
  fd: number;
  target: string;
  type: string;
};

type NetworkConnection = {
  inode: string;
  localAddress: string;
  localPort: number;
  remoteAddress: string;
  remotePort: number;
  state: string;
};

type Process = {
  pid: number;
  ppid: number;
  name: string;
  cmdline: string;
  state: string;
  threads: number;
  memoryKB: number;
  cpuTime: number;
  cpuPercent: number;
  openFiles: FileDescriptor[];
  connections: NetworkConnection[];
};

type TimelineEvent = {
  timestamp: number;
  type: string;
  pid: number;
  process: string;
  importance: "LOW" | "NORMAL" | "HIGH";
  details: Record<string, string>;
};

function getLabel(process: Process) {
  const label = process.cmdline || process.name;
  return label.length > 30 ? label.slice(0, 30) + "..." : label;
}

function getNodeColor(process: Process) {
  if (process.cpuPercent >= 50) return "#f87171";
  if (process.cpuPercent > 0) return "#fbbf24";
  if (process.state.startsWith("R")) return "#34d399";
  return "#38bdf8";
}

function getResourceIcon(type: string) {
  switch (type) {
    case "file":
      return "📄";
    case "socket":
      return "🔌";
    case "pipe":
      return "〰️";
    case "device":
      return "🖥️";
    case "kernel":
      return "⚙️";
    case "proc":
      return "🧬";
    default:
      return "❔";
  }
}

function buildTreeLayout(processes: Process[]) {
  const processMap = new Map(processes.map((p) => [p.pid, p]));
  const childrenMap = new Map<number, Process[]>();

  for (const process of processes) {
    if (!childrenMap.has(process.ppid)) {
      childrenMap.set(process.ppid, []);
    }
    childrenMap.get(process.ppid)!.push(process);
  }

  let y = 0;
  const positioned = new Map<number, { x: number; y: number }>();

  function place(pid: number, depth: number) {
    const children = childrenMap.get(pid) || [];

    positioned.set(pid, {
      x: depth * 420,
      y: y * 120,
    });

    y++;

    for (const child of children) {
      if (processMap.has(child.pid)) {
        place(child.pid, depth + 1);
      }
    }
  }

  const root = processMap.has(1) ? 1 : processes[0]?.pid;
  if (root) place(root, 0);

  for (const process of processes) {
    if (!positioned.has(process.pid)) {
      positioned.set(process.pid, { x: 0, y: y * 120 });
      y++;
    }
  }

  return positioned;
}

function App() {
  const [processes, setProcesses] = useState<Process[]>([]);
  const [events, setEvents] = useState<TimelineEvent[]>([]);
  const [search, setSearch] = useState("");
  const [selectedProcess, setSelectedProcess] = useState<Process | null>(null);
  const [eventFilter, setEventFilter] =
    useState<"ALL" | "HIGH" | "NORMAL" | "LOW">("ALL");

  useEffect(() => {
    let cancelled = false;
    let isFetchingProcesses = false;
    let isFetchingTimeline = false;

    async function fetchProcesses() {
      if (isFetchingProcesses) return;

      try {
        isFetchingProcesses = true;
        const response = await fetch("http://localhost:8080/api/processes");
        if (!response.ok) throw new Error("Failed to fetch processes");

        const data = await response.json();
        if (!cancelled) setProcesses(data);
      } catch (error) {
        console.error("KernelScope process fetch error:", error);
      } finally {
        isFetchingProcesses = false;
      }
    }

    async function fetchTimeline() {
      if (isFetchingTimeline) return;

      try {
        isFetchingTimeline = true;
        const response = await fetch("http://localhost:8080/api/timeline");
        if (!response.ok) throw new Error("Failed to fetch timeline");

        const data = await response.json();
        if (!cancelled) setEvents([...data].reverse());
      } catch (error) {
        console.error("KernelScope timeline fetch error:", error);
      } finally {
        isFetchingTimeline = false;
      }
    }

    fetchProcesses();
    fetchTimeline();

    const processInterval = setInterval(fetchProcesses, 2000);
    const timelineInterval = setInterval(fetchTimeline, 3000);

    return () => {
      cancelled = true;
      clearInterval(processInterval);
      clearInterval(timelineInterval);
    };
  }, []);

  const filteredProcesses = useMemo(() => {
    const q = search.toLowerCase().trim();
    if (!q) return processes;

    return processes.filter((p) => {
      const text = `${p.pid} ${p.ppid} ${p.name} ${p.cmdline}`.toLowerCase();
      return text.includes(q);
    });
  }, [processes, search]);

  const topCPU = useMemo(() => {
    return [...processes]
      .sort((a, b) => b.cpuPercent - a.cpuPercent)
      .slice(0, 5);
  }, [processes]);

  const topRAM = useMemo(() => {
    return [...processes]
      .sort((a, b) => b.memoryKB - a.memoryKB)
      .slice(0, 5);
  }, [processes]);

  const visibleEvents = useMemo(() => {
    if (eventFilter === "ALL") return events.slice(0, 40);

    return events
      .filter((event) => event.importance === eventFilter)
      .slice(0, 40);
  }, [events, eventFilter]);

  const { nodes, edges } = useMemo(() => {
    const visibleProcesses = filteredProcesses.slice(0, 90);
    const layout = buildTreeLayout(visibleProcesses);

    const nodes: Node[] = visibleProcesses.map((process, index) => {
      const color = getNodeColor(process);
      const isSelected = selectedProcess?.pid === process.pid;

      return {
        id: String(process.pid),
        position: layout.get(process.pid) || { x: 0, y: index * 140 },
        data: {
          label: `${getLabel(process)}
PID ${process.pid}  │  PPID ${process.ppid}
RAM ${Math.round(process.memoryKB / 1024)} MB
CPU ${process.cpuPercent.toFixed(2)}%  │  TH ${process.threads}`,
        },
        style: {
          background: isSelected ? "#1e293b" : "#0f172a",
          color: isSelected ? "#fff" : "#cbd5e1",
          border: `1.5px solid ${
            isSelected ? "#ffb04d" : "rgba(255,255,255,0.08)"
          }`,
          borderRadius: 12,
          padding: 16,
          width: 280,
          fontSize: 12,
          fontFamily: "Inter, system-ui, sans-serif",
          fontWeight: 600,
          lineHeight: 1.6,
          whiteSpace: "pre-line",
          textAlign: "left",
          boxShadow: isSelected
            ? "0 10px 35px rgba(255, 176, 77, 0.28)"
            : `0 0 18px ${color}22`,
        },
      };
    });

    const pidSet = new Set(visibleProcesses.map((p) => p.pid));

    const edges: Edge[] = visibleProcesses
      .filter((process) => pidSet.has(process.ppid))
      .map((process) => {
        const active = process.cpuPercent > 0;

        return {
          id: `${process.ppid}-${process.pid}`,
          source: String(process.ppid),
          target: String(process.pid),
          animated: active,
          type: "smoothstep",
          style: {
            stroke: active ? "#ffb04d" : "rgba(255, 255, 255, 0.12)",
            strokeWidth: active ? 2 : 1.2,
          },
        };
      });

    return { nodes, edges };
  }, [filteredProcesses, selectedProcess]);

  return (
    <main className="app">
      <section className="hero">
        <nav className="nav">
          <div>
            <h1>KERNELSCOPE</h1>
            <p>Live Linux process topology</p>
          </div>

          <input
            className="search"
            placeholder="Search PID, process, command..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />

          <div className="heroStats">
            <div>
              <strong>{processes.length}</strong>
              <span>Processes</span>
            </div>
            <div>
              <strong>
                {Math.round(
                  processes.reduce((sum, p) => sum + p.memoryKB, 0) / 1024
                )}
              </strong>
              <span>MB RAM</span>
            </div>
          </div>
        </nav>

        <div className="heroText">
          <span>LIVE KERNEL OBSERVATORY</span>
          <h2>Your operating system, visualized like a living machine.</h2>
        </div>

        <section className="graphFull">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            fitView
            defaultEdgeOptions={{ type: "smoothstep" }}
            onNodeClick={(_, node) => {
              const process = processes.find((p) => String(p.pid) === node.id);
              if (process) setSelectedProcess(process);
            }}
          >
            <Background gap={32} size={1} color="#29465f" />
            <Controls />
          </ReactFlow>
        </section>
      </section>

      <section className="infoSection">
        <div className="sectionHeader">
          <span>01</span>
          <h2>Memory Pressure</h2>
        </div>

        <div className="cardsGrid">
          {topRAM.map((p) => (
            <div className="glassCard" key={p.pid}>
              <small>PID {p.pid}</small>
              <h3>{getLabel(p)}</h3>
              <strong>{Math.round(p.memoryKB / 1024)} MB</strong>
            </div>
          ))}
        </div>
      </section>

      <section className="infoSection darkBand">
        <div className="sectionHeader">
          <span>02</span>
          <h2>CPU Activity</h2>
        </div>

        <div className="cardsGrid">
          {topCPU.map((p) => (
            <div className="glassCard hotCard" key={p.pid}>
              <small>PID {p.pid}</small>
              <h3>{getLabel(p)}</h3>
              <strong>{p.cpuPercent.toFixed(2)}%</strong>
            </div>
          ))}
        </div>
      </section>

      <section className="infoSection timelineSection">
        <div className="sectionHeader">
          <span>03</span>
          <h2>Process Timeline</h2>
          <p>
            A live chronological stream of process starts, exits, socket
            activity, resource changes, and kernel-level signals.
          </p>
        </div>

        <div className="timelineFilters">
          {(["ALL", "HIGH", "NORMAL", "LOW"] as const).map((filter) => (
            <button
              key={filter}
              className={eventFilter === filter ? "activeFilter" : ""}
              onClick={() => setEventFilter(filter)}
            >
              {filter}
            </button>
          ))}
        </div>

        <div className="timelineTrack">
          {visibleEvents.map((event, index) => (
            <article
              className={`timelineCard ${event.importance.toLowerCase()}`}
              key={`${event.timestamp}-${event.pid}-${index}`}
            >
              <div className="timelineMarker" />

              <div className="timelineContent">
                <div className="timelineTop">
                  <span>{event.type}</span>
                  <b>{event.importance}</b>
                </div>

                <h3>{event.process}</h3>

                <div className="timelineMeta">
                  <p>PID {event.pid}</p>
                  <p>{new Date(event.timestamp).toLocaleTimeString()}</p>
                </div>

                {event.details && Object.keys(event.details).length > 0 && (
                  <div className="timelineDetails">
                    {Object.entries(event.details)
                      .slice(0, 3)
                      .map(([key, value]) => (
                        <p key={key}>
                          <b>{key}</b>
                          <span>{String(value)}</span>
                        </p>
                      ))}
                  </div>
                )}
              </div>
            </article>
          ))}
        </div>
      </section>

      <section className="infoSection">
        <div className="sectionHeader">
          <span>04</span>
          <h2>Process Inspector</h2>
        </div>

        {!selectedProcess ? (
          <div className="emptyInspector">
            Click any node in the graph above to inspect files, sockets and live
            network connections.
          </div>
        ) : (
          <div className="inspector">
            <div className="processMain">
              <h3>{getLabel(selectedProcess)}</h3>
              <p>
                <b>PID:</b> {selectedProcess.pid}
              </p>
              <p>
                <b>PPID:</b> {selectedProcess.ppid}
              </p>
              <p>
                <b>State:</b> {selectedProcess.state}
              </p>
              <p>
                <b>CPU:</b> {selectedProcess.cpuPercent.toFixed(2)}%
              </p>
              <p>
                <b>RAM:</b> {Math.round(selectedProcess.memoryKB / 1024)} MB
              </p>
              <p>
                <b>Threads:</b> {selectedProcess.threads}
              </p>
            </div>

            <div>
              <h3>
                Network Connections ({selectedProcess.connections?.length || 0})
              </h3>

              <div className="listStack">
                {(selectedProcess.connections || []).map((connection) => (
                  <div className="connectionRow" key={connection.inode}>
                    <b>🌐 {connection.state}</b>
                    <p>
                      {connection.localAddress}:{connection.localPort}
                    </p>
                    {connection.state !== "LISTEN" && (
                      <p>
                        → {connection.remoteAddress}:{connection.remotePort}
                      </p>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <div>
              <h3>Open Resources ({selectedProcess.openFiles?.length || 0})</h3>

              <div className="listStack">
                {(selectedProcess.openFiles || []).slice(0, 70).map((file) => (
                  <div
                    className="resourceRow"
                    key={`${selectedProcess.pid}-${file.fd}`}
                  >
                    <span>{getResourceIcon(file.type)}</span>
                    <div>
                      <b>FD {file.fd}</b>
                      <p>{file.target}</p>
                      <small>{file.type}</small>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </section>
    </main>
  );
}

export default App;