# 📡 kafka-spark-realtime — Real-Time Clickstream Analytics Pipeline

> A complete end-to-end pipeline for real-time analytics using Kafka + Spark Structured Streaming + PostgreSQL + Elasticsearch.

---

## 📌 Overview

It simulates and processes clickstream data in real time. It uses:

- **Kafka** to ingest user click events
- **Spark Structured Streaming** (via PySpark) for stream processing
- **PostgreSQL** to store historical data
- **Elasticsearch + Kibana** for search & visualization

Optional tools like Kafka Connect and Kafdrop are also included.

---

## 🧱 Architecture

```

\[Simulator (Go)]
↓
\[Kafka Topic: clicks]
↓
\[Spark Structured Streaming]
↓           ↘
\[PostgreSQL]   \[Elasticsearch]
↓
\[Kibana]

```

---

## 🚀 Features

- ✅ Real-time event simulation (Go)
- ✅ Kafka producer/consumer pipeline
- ✅ Stream processing with PySpark
- ✅ Sink to PostgreSQL and Elasticsearch
- ✅ Docker Compose setup for full reproducibility
- ✅ Extensible and production-oriented

---

## 🛠️ Tech Stack

| Tool              | Role                                |
|-------------------|-------------------------------------|
| **Kafka**         | Message broker for real-time events |
| **Spark**         | Stream processing engine (PySpark)  |
| **PostgreSQL**    | Storage of historical events        |
| **Elasticsearch** | Real-time search and analytics      |
| **Kibana**        | Visualization UI (optional)         |
| **Kafka Connect** | Sink connectors (optional)          |
| **Docker Compose**| Environment orchestration           |

---

## 📂 Project Structure

```

kafka-spark-realtime/
├── docker-compose.yml
├── simulator/
│   └── main.go
├── spark/
│   └── stream\_processor.py
├── sql/
│   └── create\_tables.sql
├── connectors/
│   ├── postgres-sink.json
│   └── elastic-sink.json
├── notebooks/
│   └── exploratory\_analysis.ipynb
└── README.md

````

---

## ▶️ Quick Start

1. **Clone the repo**
```bash
git clone [https://github.com//kafka-spark-realtime.git](https://github.com/Yurhigz/kafka-spark-realtime.git)
cd kafka-spark-realtime
````

2. **Launch the environment**

```bash
docker-compose up -d
```

2.bis **Setup C Bindings**

Installer le compilateur GCC : 
```bash
sudo apt-get update
sudo apt-get install build-essential
```

Modifier la variable CGO_ENABLED

3. **Start the simulator**

```bash
go run simulator/main.go

```

4. **Launch the Spark job**

```bash
spark-submit spark/stream_processor.py
```

5. **Explore**

* PostgreSQL: `localhost:5432`
* Elasticsearch: `localhost:9200`
* Kibana: `localhost:5601`
* Kafdrop: `localhost:9000`

---

## 📊 Example Event Format

```json
{
  "user_id": "user_42",
  "page": "/product/1234",
  "event": "click",
  "timestamp": "2025-06-04T12:00:00Z"
}
```

---

## 🤝 License

MIT License — feel free to use for academic or enterprise learning.

---
