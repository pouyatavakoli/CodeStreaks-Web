# CodeStreaks Web
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)

CodeStreaks Web is a **Go-based web application** that provides the same core functionality as the original **CodeStreaks** project, but exposed through a **web interface and HTTP API** instead of a Python CLI.

This project is a **full reimplementation in Go**, designed for better performance, easier deployment, and web-based usage.


## Background

The original project, **CodeStreaks**, was implemented in Python as a command-line tool for analyzing Codeforces submission data and computing daily streaks and leaderboards.

* Original repository (Python CLI):
  **[CodeStreaks](https://github.com/pouyatavakoli/CodeStreaks)**

This repository (**CodeStreaks-Web**) is **not a fork**. It is a separate codebase that:

* Reimplements the same logic in **Go**
* Targets a **web/server architecture**

The Python project serves as the **reference specification** for behavior and expected results.

---

## Relationship to CodeStreaks (Python)

| Aspect    | CodeStreaks (Python) | CodeStreaks Web (Go) |
| --------- | -------------------- | -------------------- |
| Language  | Python               | Go                   |
| Interface | CLI                  | Web + HTTP API       |
| Purpose   | Script / automation  | Service / website    |
| Status    | Stable reference     | Active development   |
