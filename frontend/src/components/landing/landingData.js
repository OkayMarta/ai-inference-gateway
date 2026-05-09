export const landingFeatureItems = [
    {
        icon: "prompt",
        tone: "blue",
        title: "Prompt orchestration",
        description: "Submit requests through one unified gateway.",
    },
    {
        icon: "billing",
        tone: "violet",
        title: "Token billing",
        description: "Track balance, pricing, and automatic deductions.",
    },
    {
        icon: "tasks",
        tone: "blue",
        title: "Task monitoring",
        description: "See statuses like Queued, Processing, Completed, Failed, and Cancelled.",
    },
    {
        icon: "history",
        tone: "green",
        title: "Result history",
        description: "Open previous tasks and review generated outputs anytime.",
    },
];

export const landingModelItems = [
    {
        name: "qwen2.5:1.5b",
        description: "Balanced model for everyday prompts, summaries, and general reasoning.",
        tags: ["General use", "Summaries", "Reasoning"],
        recommended: true,
    },
    {
        name: "tinyllama:latest",
        description: "Lightweight model for quick simple tasks and fast responses.",
        tags: ["Quick tasks", "Simple Q&A", "Fast responses"],
        recommended: false,
    },
];

export const landingPricingPlans = [
    {
        name: "Starter",
        price: "$10",
        units: "1,000 units",
        description: "Good for trying the platform",
        cta: "Choose Starter",
        recommended: false,
    },
    {
        name: "Standard",
        price: "$45",
        units: "5,000 units",
        description: "Best value for regular usage",
        cta: "Choose Standard",
        recommended: true,
    },
    {
        name: "Pro",
        price: "$99",
        units: "12,000 units",
        description: "For heavier workloads and teams",
        cta: "Choose Pro",
        recommended: false,
    },
];
