import "./FilterTabs.css";

interface FilterTab<T extends string> {
        label: string;
        value: T;
        count?: number;
}

interface Props<T extends string> {
        tabs: FilterTab<T>[];
        active: T;
        onChange: (value: T) => void;
}

export default function FilterTabs<T extends string>({ tabs, active, onChange }: Props<T>) {
        return (
                <div className="filter-tabs">
                        {tabs.map((tab) => (
                                <button
                                        key={tab.value}
                                        className={`filter-tab ${active === tab.value ? "filter-tab--active" : ""}`}
                                        onClick={() => onChange(tab.value)}
                                >
                                        {tab.label}
                                        {tab.count !== undefined && (
                                                <span className="filter-tab-count">
                                                        {tab.count}
                                                </span>
                                        )}
                                </button>
                        ))}
                </div>
        );
}
