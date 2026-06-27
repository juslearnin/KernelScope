import { useEffect, useMemo, useState } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
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
};

function getLabel(process: Process) {
  const label = process.cmdline || process.name;
  return label.length > 42 ? label.slice(0, 42) + "..." : label;
}

function getNodeColor(process: Process) {
  if (process.cpuPercent >= 50) return "#ff4d6d";
  if (process.cpuPercent > 0) return "#ffd166";
  if (process.state.startsWith("R")) return "#4ade80";
  return "#7df9ff";
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

    const currentY = y;
    positioned.set(pid, {
      x: depth * 300,
      y: currentY * 120,
    });

    y++;

    for (const child of children) {
      if (processMap.has(child.pid)) {
        place(child.pid, depth + 1);
      }
    }
  }

  const root = processMap.has(1) ? 1 : processes[0]?.pid;

  if (root) {
    place(root, 0);
  }

  for (const process of processes) {
    if (!positioned.has(process.pid)) {
      positioned.set(process.pid, {
        x: 0,
        y: y * 120,
      });
      y++;
    }
  }

  return positioned;
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

function App() {
  const [processes, setProcesses] = useState<Process[]>([]);
  const [search, setSearch] = useState("");
  const [selectedProcess, setSelectedProcess] = useState<Process | null>(null);

  useEffect(() => {
    let cancelled = false;
    let isFetching = false;

    async function fetchProcesses() {
      if (isFetching) return;

      try {
        isFetching = true;
        const response = await fetch("http://localhost:8080/api/processes");

        if (!response.ok) throw new Error("Failed to fetch processes");

        const data = await response.json();

        if (!cancelled) setProcesses(data);
      } catch (error) {
        console.error("KernelScope fetch error:", error);
      } finally {
        isFetching = false;
      }
    }

    fetchProcesses();
    const interval = setInterval(fetchProcesses, 2000);

    return () => {
      cancelled = true;
      clearInterval(interval);
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

  const { nodes, edges } = useMemo(() => {
    const visibleProcesses = filteredProcesses.slice(0, 90);
    const layout = buildTreeLayout(visibleProcesses);

    const nodes: Node[] = visibleProcesses.map((process, index) => {
      const color = getNodeColor(process);

      return {
        id: String(process.pid),
        position: layout.get(process.pid) || {
          x: 0,
          y: index * 120,
        },
        data: {
          label: `${getLabel(process)}
PID: ${process.pid}
RAM: ${process.memoryKB} KB
CPU: ${process.cpuPercent.toFixed(2)}
Threads: ${process.threads}`,
        },
        style: {
          background: "#0b1220",
          color,
          border: `1px solid ${color}`,
          borderRadius: 14,
          padding: 12,
          width: 230,
          fontSize: 12,
          whiteSpace: "pre-line",
          boxShadow: `0 0 18px ${color}55`,
        },
      };
    });

    const pidSet = new Set(visibleProcesses.map((p) => p.pid));

    const edges: Edge[] = visibleProcesses
      .filter((process) => pidSet.has(process.ppid))
      .map((process) => ({
        id: `${process.ppid}-${process.pid}`,
        source: String(process.ppid),
        target: String(process.pid),
        animated: process.cpuPercent > 0,
        style: {
          stroke: process.cpuPercent > 0 ? "#ffd166" : "#355070",
        },
      }));

    return { nodes, edges };
  }, [filteredProcesses]);

  return (
    <main className="app">
      <aside className="sidebar">
        <h1>KERNELSCOPE</h1>
        <p className="subtitle">Live OS Process Map</p>

        <input
          className="search"
          placeholder="Search pid, node, bash..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />

        <div className="statGrid">
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

        <section>
          <h2>Top CPU</h2>
          {topCPU.map((p) => (
            <div className="miniRow" key={p.pid}>
              <span>{getLabel(p)}</span>
              <b>{p.cpuPercent.toFixed(2)}</b>
            </div>
          ))}
        </section>

        <section>
          <h2>Top RAM</h2>
          {topRAM.map((p) => (
            <div className="miniRow" key={p.pid}>
              <span>{getLabel(p)}</span>
              <b>{Math.round(p.memoryKB / 1024)} MB</b>
            </div>
          ))}
        </section>

        <section>
          <h2>Legend</h2>
          <p><span className="dot idle" /> Sleeping / waiting</p>
          <p><span className="dot running" /> Running</p>
          <p><span className="dot active" /> CPU active</p>
          <p><span className="dot hot" /> CPU heavy</p>
        </section>
      </aside>

      <section className="graph">
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
          <Background />
          <Controls />
          <MiniMap />
        </ReactFlow>
      </section>

      <aside className="detailsPanel">
        {!selectedProcess ? (
          <div className="emptyDetails">
            <h2>Process Details</h2>
            <p>Click a process node to inspect its open resources.</p>
          </div>
        ) : (
          <>
            <h2>{getLabel(selectedProcess)}</h2>

            <div className="detailBlock">
              <p><b>PID:</b> {selectedProcess.pid}</p>
              <p><b>PPID:</b> {selectedProcess.ppid}</p>
              <p><b>State:</b> {selectedProcess.state}</p>
              <p><b>CPU:</b> {selectedProcess.cpuPercent.toFixed(2)}</p>
              <p><b>RAM:</b> {Math.round(selectedProcess.memoryKB / 1024)} MB</p>
              <p><b>Threads:</b> {selectedProcess.threads}</p>
            </div>

            <h3>Open Resources ({selectedProcess.openFiles?.length || 0})</h3>

            <div className="resources">
              {(selectedProcess.openFiles || []).slice(0, 80).map((file) => (
                <div className="resourceRow" key={`${selectedProcess.pid}-${file.fd}`}>
                  <span className="resourceIcon">{getResourceIcon(file.type)}</span>
                  <div>
                    <b>FD {file.fd}</b>
                    <p>{file.target}</p>
                    <small>{file.type}</small>
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
      </aside>
    </main>
  );
}

export default App;