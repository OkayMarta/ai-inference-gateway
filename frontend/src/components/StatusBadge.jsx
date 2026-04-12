import "./StatusBadge.css";

function getStatusClass(status) {
    switch (status) {
        case "Completed":
            return "status-badge status-completed";
        case "Processing":
            return "status-badge status-processing";
        case "Failed":
            return "status-badge status-failed";
        default:
            return "status-badge status-queued";
    }
}

export default function StatusBadge({ status }) {
    return <span className={getStatusClass(status)}>{status}</span>;
}
