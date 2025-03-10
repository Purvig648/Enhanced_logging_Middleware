# Enhanced Logging Middleware

## 📌 Overview
`Enhanced Logging Middleware` is a Golang library that provides structured logging using **Logrus** with features like:
- **Log rotation** (via Lumberjack)
- **Middleware logging for HTTP requests**
- **Async logging for better performance**
- **External log sinks** (e.g., Loki, CloudWatch, Elasticsearch)
- **Trace ID generation** for request tracking

## 🚀 Features
- 🔹 **JSON/Text Logging**
- 🔹 **Customizable Log Levels**
- 🔹 **Concurrency-Safe Logging**
- 🔹 **Middleware for HTTP Logging**
- 🔹 **Pluggable Logging Sinks**
- 🔹 **Trace ID Propagation**
- 🔹 **Asynchronous Logging Support**

---

## 🛠 Installation
To install the package, run:
```sh
go get github.com/Purvig648/Enhanced_logging_Middleware
```

---

## 📂 Project Structure
```
Enhanced_logging_Middleware/
│── logger/
│   ├── logger.go        # Core logging functionality
│── middleware/
│   ├── middleware.go    # HTTP Middleware logic
│── go.mod              # Module file
│── README.md           # Project documentation
│── main.go             # Example usage
```

---

## 🤝 Contributing
Feel free to open issues or submit pull requests to enhance functionality!

## 📧 Contact
For any queries, reach out to **Purvig648** via GitHub.

