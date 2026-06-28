# KernelScope

### Real-Time Linux Process Observability Platform

> **Think of it as "Google Maps for your Linux Operating System."**
>
> KernelScope is a real-time Linux observability platform that visualizes processes, system resources, open files, network connections, and process lifecycles by directly parsing the Linux `/proc` filesystem. It combines low-level systems programming with an interactive React-based visualization layer to help developers understand how an operating system behaves in real time.

---

<img width="851" height="442" alt="image" src="https://github.com/user-attachments/assets/c64b0ca6-1fd0-41aa-84a0-9924810480d3" />



<img width="812" height="338" alt="image" src="https://github.com/user-attachments/assets/ddbaec6c-4917-4c37-9594-3a0c668d18ed" />



<img width="802" height="404" alt="image" src="https://github.com/user-attachments/assets/f29f6c07-13cc-4520-bb78-9daf217d4f3c" />



<img width="941" height="405" alt="image" src="https://github.com/user-attachments/assets/57792ae5-2115-408d-8592-2346b32912b3" />



<img width="923" height="380" alt="image" src="https://github.com/user-attachments/assets/54eb6207-ae48-4fd4-b75b-c8b8c5b272c6" />




# Need of KernelScope?

Modern operating systems execute hundreds of processes simultaneously.

Traditional tools like `ps`, `top`, `htop`, `lsof`, `ss`, and `netstat` provide valuable information, but each focuses on only one aspect of the system.

KernelScope combines these views into a single interactive platform that answers questions such as:

* Which process created this child process?
* Which files does this process currently have open?
* Which TCP connections belong to this application?
* What changed in the system during the last minute?
* Which processes triggered alerts?
* How are all processes related?

Instead of reading raw terminal output, KernelScope presents the operating system as an interconnected live graph.

---

# Features

### Process Observatory

* Real-time process discovery
* Parent-child process hierarchy
* Process tree visualization
* Command line extraction
* Process state monitoring
* Thread count
* Memory usage
* CPU usage estimation

---

### Resource Observatory

* Live CPU utilization
* Resident memory tracking
* Thread monitoring
* Process statistics

---

### File Observatory

* Open file descriptor discovery
* File descriptor classification
* File path resolution
* Socket identification
* Anonymous inode detection
* Device identification

---

### Network Observatory

* TCP socket parsing
* Socket-to-process mapping
* Local and remote address extraction
* Connection state monitoring
* Listening socket discovery

---

### Timeline Engine

KernelScope continuously compares process snapshots to detect lifecycle events.

Supported events include:

* Process Started
* Process Exited

Each event is automatically timestamped and stored.

---

### Rule Engine

KernelScope evaluates live process data against configurable rules.

Current rules include:

* High CPU
* High Memory
* Many Open Files
* High Network Connections

Alerts are automatically deduplicated to prevent repeated notifications.

---

### Alert Manager

KernelScope tracks active alerts separately from rule evaluation.

Features:

* Duplicate suppression
* Active alert tracking
* Alert history
* REST API exposure

---

### SQLite Persistence

Timeline events are persisted using SQLite.

Benefits:

* Survives backend restarts
* Zero external database dependencies
* Lightweight embedded storage

---

### Real-Time Streaming

KernelScope uses WebSockets to stream events directly to connected clients.

Unlike polling, updates are pushed immediately when the operating system changes.

---

# System Architecture

```text
                 Linux Kernel
                      │
                      ▼
               /proc Filesystem
                      │
      ┌───────────────┴───────────────┐
      ▼                               ▼
 Process Collector             Network Collector
      ▼                               ▼
 Resource Collector            File Collector
      └───────────────┬───────────────┘
                      ▼
               Process Graph Builder
                      ▼
               Timeline Engine
                      ▼
                Rule Evaluation
                      ▼
                 Alert Manager
               ┌──────────────┐
               ▼              ▼
          SQLite Storage   WebSocket Hub
               └──────────────┘
                      ▼
                React Frontend
```

---

# Technology Stack

## Backend

* Go
* SQLite
* Gorilla WebSocket
* REST API

## Frontend

* React
* TypeScript
* React Flow

## Operating System

* Linux
* `/proc` Filesystem

---

# Internal Modules

```
backend/

collector/
timeline/
rules/
storage/
server/
websocket/
models/
```

Each module has a single responsibility, making the architecture modular and extensible.

---

# Learning Outcomes

KernelScope was built as a deep dive into Linux internals and observability systems.

The project explores concepts such as:

* Linux process management
* `/proc` filesystem
* Process scheduling
* Parent-child relationships
* File descriptors
* TCP networking
* Process lifecycle detection
* Event-driven architecture
* Rule engines
* SQLite persistence
* WebSocket communication
* Real-time visualization

---

# Future Improvements

The current implementation focuses on Linux process observability.

Potential future enhancements include:

* Session recording
* Replay engine
* eBPF integration
* Container awareness
* Memory map visualization
* Disk I/O monitoring
* Environment variable inspection
* Plugin-based collectors
* Export functionality
* Advanced rule definitions

---

# Project Vision

KernelScope aims to make Linux internals intuitive and interactive.

Rather than treating the operating system as a collection of disconnected command-line utilities, KernelScope visualizes it as a dynamic, living system where processes, files, sockets, and resources continuously interact.

The long-term vision is to become an educational and debugging platform that helps developers understand operating system behavior through real-time observability.

---

# Why KernelScope?

Modern operating systems expose an enormous amount of runtime information through the Linux `/proc` filesystem.

Although individual tools exist to inspect processes, sockets, files, and resources, they often provide isolated views of the system.

KernelScope brings these perspectives together into a single platform capable of visualizing relationships between processes, files, sockets, resources, timelines, and alerts.

The goal is not only to monitor a Linux system, but also to make its behavior easier to understand.

# Acknowledgements

KernelScope was developed as a systems programming project to explore Linux internals, observability concepts, and modern backend architecture using Go and React.

Every feature—from parsing `/proc` to real-time WebSocket streaming—was implemented to deepen understanding of how Linux exposes kernel information and how monitoring platforms are designed.

# What This Project Demonstrates

KernelScope demonstrates practical experience with:

- Systems Programming
- Linux Internals
- Operating Systems
- Go Backend Development
- Concurrent Programming
- REST APIs
- WebSockets
- SQLite
- React
- Event-Driven Systems
- Rule Engines
- Observability
- Software Architecture
