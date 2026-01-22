#!/usr/bin/env python3

import re
import sys
import matplotlib.pyplot as plt

def parse_tcmalloc_stats(filename):
    classes = []
    pattern = re.compile(
        r"class\s+(\d+)\s+\[\s+(\d+)\s+bytes\s+\]\s*:\s*"
        r"(\d+)\s+objs;\s+([\d\.]+)\s+MiB;.*?"
        r"spans:\s+(\d+)\s+ret\s*/\s*(\d+)\s+req\s*=\s*([\d\.]+)"
    )
    with open(filename, "r") as f:
        for line in f:
            m = pattern.search(line)
            if m:
                class_id = int(m.group(1))
                obj_size = int(m.group(2))
                objs = int(m.group(3))
                mbytes = float(m.group(4))
                spans_ret = int(m.group(5))
                spans_req = int(m.group(6))
                reuse_ratio = float(m.group(7))
                classes.append({
                    "class": class_id,
                    "size": obj_size,
                    "objs": objs,
                    "mib": mbytes,
                    "spans_ret": spans_ret,
                    "spans_req": spans_req,
                    "reuse": reuse_ratio
                })
    return classes

def plot_top_classes(classes, topn=5):
    # Sort by memory usage
    top = sorted(classes, key=lambda c: c["mib"], reverse=True)[:topn]

    sizes = [c["size"] for c in top]
    mib = [c["mib"] for c in top]
    reuse = [c["reuse"] for c in top]

    # Highlight classes with poor reuse ratio
    colors = ["red" if r < 0.5 else "skyblue" for r in reuse]

    fig, ax1 = plt.subplots(figsize=(8,5))

    # Bar chart for memory usage
    ax1.bar(range(len(top)), mib, color=colors, alpha=0.7, label="Memory (MiB)")
    ax1.set_ylabel("Memory (MiB)")
    ax1.set_xticks(range(len(top)))
    ax1.set_xticklabels([f"{s}B" for s in sizes], rotation=45)

    # Line plot for reuse ratio
    ax2 = ax1.twinx()
    ax2.plot(range(len(top)), reuse, color="darkred", marker="o", label="Reuse ratio")
    ax2.set_ylabel("Span reuse ratio")

    # Add threshold line at 0.5
    ax2.axhline(0.5, color="gray", linestyle="--", linewidth=0.8)

    # Legends
    ax1.legend(loc="upper left")
    ax2.legend(loc="upper right")
    plt.title(f"Top {topn} TCMalloc Size Classes by Memory Usage")
    plt.tight_layout()
    plt.show()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <tcmalloc-stats-file>")
        sys.exit(1)
    stats_file = sys.argv[1]
    classes = parse_tcmalloc_stats(stats_file)
    plot_top_classes(classes, topn=5)
