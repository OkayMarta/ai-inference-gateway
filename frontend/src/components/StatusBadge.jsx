import "../styles/components/StatusBadge.css";

const getStatusClass = (status) => {
    switch (status) {
        case "Completed":
            return "status-badge status-completed";
        case "Processing":
            return "status-badge status-processing";
        case "Failed":
            return "status-badge status-failed";
        case "Cancelled":
            return "status-badge status-cancelled";
        default:
            return "status-badge status-queued";
    }
};

const StatusBadge = ({ status }) => {
    return <span className={getStatusClass(status)}>{status}</span>;
};

export default StatusBadge;
