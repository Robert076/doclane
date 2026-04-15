"use client";
import { Stats } from "@/types/stats";
import {
        AreaChart,
        Area,
        BarChart,
        Bar,
        XAxis,
        YAxis,
        CartesianGrid,
        Tooltip,
        ResponsiveContainer,
} from "recharts";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { generateStatsPDF } from "@/lib/client/generateStatsPDF";
import "./StatsSection.css";

interface Props {
        stats: Stats;
}

interface StatCardProps {
        label: string;
        value: string | number;
        sub?: string;
        change?: number;
}

function StatCard({ label, value, sub, change }: StatCardProps) {
        return (
                <div className="stat-card">
                        <span className="stat-card-label">{label}</span>
                        <span className="stat-card-value">{value}</span>
                        {sub && <span className="stat-card-sub">{sub}</span>}
                        {change !== undefined && (
                                <span
                                        className={`stat-card-change ${change >= 0 ? "stat-card-change--up" : "stat-card-change--down"}`}
                                >
                                        {change >= 0 ? "▲" : "▼"} {Math.abs(change).toFixed(1)}
                                        %
                                </span>
                        )}
                </div>
        );
}

function SectionTitle({ text }: { text: string }) {
        return <h2 className="stats-section-title">{text}</h2>;
}

function formatHours(hours: number): string {
        if (hours === 0) return "N/A";
        const days = Math.floor(hours / 24);
        const remainingHours = Math.round(hours % 24);
        if (days === 0) return `${remainingHours}h`;
        if (remainingHours === 0) return `${days}z`;
        return `${days}z ${remainingHours}h`;
}

export default function StatsSection({ stats }: Props) {
        return (
                <div className="stats-section">
                        <div className="stats-export">
                                <ButtonPrimary
                                        text="Exportă PDF"
                                        variant="ghost"
                                        onClick={() => generateStatsPDF(stats)}
                                />
                        </div>

                        <SectionTitle text="Cereri" />
                        <div className="stats-grid">
                                <StatCard
                                        label="Cereri deschise"
                                        value={stats.total_open_requests}
                                />
                                <StatCard
                                        label="Cereri finalizate"
                                        value={stats.total_archived_requests}
                                />
                                <StatCard
                                        label="Cereri retrase"
                                        value={stats.total_cancelled_requests}
                                />
                                <StatCard
                                        label="Rată finalizare"
                                        value={`${stats.completion_rate.toFixed(1)}%`}
                                />
                                <StatCard
                                        label="Rată retragere"
                                        value={`${stats.cancellation_rate.toFixed(1)}%`}
                                />
                                <StatCard
                                        label="Timp mediu finalizare"
                                        value={formatHours(stats.avg_completion_hours)}
                                />
                        </div>

                        <SectionTitle text="Activitate" />
                        <div className="stats-grid">
                                <StatCard
                                        label="Cereri săptămâna aceasta"
                                        value={stats.requests_this_week}
                                        sub={`față de ${stats.requests_last_week} săptămâna trecută`}
                                        change={stats.weekly_change_percent}
                                />
                                <StatCard
                                        label="Cereri luna aceasta"
                                        value={stats.requests_this_month}
                                        sub={`față de ${stats.requests_last_month} luna trecută`}
                                        change={stats.monthly_change_percent}
                                />
                        </div>

                        <div className="stats-chart-card">
                                <h3 className="stats-chart-title">
                                        Cereri în ultimele 7 zile
                                </h3>
                                <ResponsiveContainer width="100%" height={220}>
                                        <AreaChart data={stats.requests_last_7_days}>
                                                <defs>
                                                        <linearGradient
                                                                id="colorCount"
                                                                x1="0"
                                                                y1="0"
                                                                x2="0"
                                                                y2="1"
                                                        >
                                                                <stop
                                                                        offset="5%"
                                                                        stopColor="#FF5722"
                                                                        stopOpacity={0.15}
                                                                />
                                                                <stop
                                                                        offset="95%"
                                                                        stopColor="#FF5722"
                                                                        stopOpacity={0}
                                                                />
                                                        </linearGradient>
                                                </defs>
                                                <CartesianGrid
                                                        strokeDasharray="3 3"
                                                        stroke="#f1f5f9"
                                                />
                                                <XAxis
                                                        dataKey="date"
                                                        tick={{
                                                                fontSize: 12,
                                                                fill: "#94a3b8",
                                                        }}
                                                />
                                                <YAxis
                                                        allowDecimals={false}
                                                        tick={{
                                                                fontSize: 12,
                                                                fill: "#94a3b8",
                                                        }}
                                                />
                                                <Tooltip />
                                                <Area
                                                        type="monotone"
                                                        dataKey="count"
                                                        stroke="#FF5722"
                                                        strokeWidth={2}
                                                        fill="url(#colorCount)"
                                                />
                                        </AreaChart>
                                </ResponsiveContainer>
                        </div>

                        <SectionTitle text="Departamente" />
                        <div className="stats-grid">
                                <StatCard
                                        label="Total departamente"
                                        value={stats.total_departments}
                                />
                                <StatCard
                                        label="Membri departamente"
                                        value={stats.total_department_members}
                                />
                        </div>

                        <div className="stats-chart-card">
                                <h3 className="stats-chart-title">
                                        Cereri deschise per departament
                                </h3>
                                <ResponsiveContainer width="100%" height={220}>
                                        <BarChart data={stats.requests_per_department}>
                                                <CartesianGrid
                                                        strokeDasharray="3 3"
                                                        stroke="#f1f5f9"
                                                />
                                                <XAxis
                                                        dataKey="department_name"
                                                        tick={{
                                                                fontSize: 12,
                                                                fill: "#94a3b8",
                                                        }}
                                                />
                                                <YAxis
                                                        allowDecimals={false}
                                                        tick={{
                                                                fontSize: 12,
                                                                fill: "#94a3b8",
                                                        }}
                                                />
                                                <Tooltip />
                                                <Bar
                                                        dataKey="request_count"
                                                        fill="#FF5722"
                                                        radius={[4, 4, 0, 0]}
                                                />
                                        </BarChart>
                                </ResponsiveContainer>
                        </div>

                        {stats.requests_per_locality?.length > 0 && (
                                <div className="stats-chart-card">
                                        <h3 className="stats-chart-title">
                                                Cereri per localitate (top 10)
                                        </h3>
                                        <ResponsiveContainer width="100%" height={280}>
                                                <BarChart
                                                        data={stats.requests_per_locality}
                                                        layout="vertical"
                                                >
                                                        <CartesianGrid
                                                                strokeDasharray="3 3"
                                                                stroke="#f1f5f9"
                                                        />
                                                        <XAxis
                                                                type="number"
                                                                allowDecimals={false}
                                                                tick={{
                                                                        fontSize: 12,
                                                                        fill: "#94a3b8",
                                                                }}
                                                        />
                                                        <YAxis
                                                                dataKey="locality"
                                                                type="category"
                                                                tick={{
                                                                        fontSize: 12,
                                                                        fill: "#94a3b8",
                                                                }}
                                                                width={140}
                                                        />
                                                        <Tooltip />
                                                        <Bar
                                                                dataKey="request_count"
                                                                fill="#FF5722"
                                                                radius={[0, 4, 4, 0]}
                                                        />
                                                </BarChart>
                                        </ResponsiveContainer>
                                </div>
                        )}

                        <SectionTitle text="Utilizatori" />
                        <div className="stats-grid">
                                <StatCard
                                        label="Total utilizatori"
                                        value={stats.total_users}
                                />
                                <StatCard label="Cetățeni" value={stats.total_citizens} />
                                <StatCard
                                        label="Utilizatori activi"
                                        value={stats.total_active_users}
                                />
                                <StatCard
                                        label="Utilizatori dezactivați"
                                        value={stats.total_deactivated_users}
                                />
                        </div>

                        <SectionTitle text="Șabloane" />
                        <div className="stats-grid">
                                <StatCard
                                        label="Șabloane active"
                                        value={stats.total_active_templates}
                                />
                                <StatCard
                                        label="Șabloane arhivate"
                                        value={stats.total_archived_templates}
                                />
                        </div>

                        <div className="stats-chart-card">
                                <h3 className="stats-chart-title">Top 5 șabloane folosite</h3>
                                <ResponsiveContainer width="100%" height={220}>
                                        <BarChart
                                                data={stats.most_used_templates}
                                                layout="vertical"
                                        >
                                                <CartesianGrid
                                                        strokeDasharray="3 3"
                                                        stroke="#f1f5f9"
                                                />
                                                <XAxis
                                                        type="number"
                                                        allowDecimals={false}
                                                        tick={{
                                                                fontSize: 12,
                                                                fill: "#94a3b8",
                                                        }}
                                                />
                                                <YAxis
                                                        dataKey="template_title"
                                                        type="category"
                                                        tick={{
                                                                fontSize: 12,
                                                                fill: "#94a3b8",
                                                        }}
                                                        width={180}
                                                />
                                                <Tooltip />
                                                <Bar
                                                        dataKey="request_count"
                                                        fill="#FF5722"
                                                        radius={[0, 4, 4, 0]}
                                                />
                                        </BarChart>
                                </ResponsiveContainer>
                        </div>

                        {stats.member_stats?.length > 0 && (
                                <>
                                        <SectionTitle text="Performanță membri" />
                                        <div className="stats-table-card">
                                                <table className="stats-table">
                                                        <thead>
                                                                <tr>
                                                                        <th>Membru</th>
                                                                        <th>Departament</th>
                                                                        <th>Preluate</th>
                                                                        <th>Finalizate</th>
                                                                        <th>În lucru</th>
                                                                        <th>Timp mediu</th>
                                                                </tr>
                                                        </thead>
                                                        <tbody>
                                                                {stats.member_stats.map(
                                                                        (m) => (
                                                                                <tr
                                                                                        key={
                                                                                                m.user_id
                                                                                        }
                                                                                >
                                                                                        <td>
                                                                                                {
                                                                                                        m.first_name
                                                                                                }{" "}
                                                                                                {
                                                                                                        m.last_name
                                                                                                }
                                                                                        </td>
                                                                                        <td>
                                                                                                {
                                                                                                        m.department_name
                                                                                                }
                                                                                        </td>
                                                                                        <td>
                                                                                                {
                                                                                                        m.total_claimed
                                                                                                }
                                                                                        </td>
                                                                                        <td>
                                                                                                {
                                                                                                        m.total_closed
                                                                                                }
                                                                                        </td>
                                                                                        <td>
                                                                                                {
                                                                                                        m.total_pending
                                                                                                }
                                                                                        </td>
                                                                                        <td>
                                                                                                {formatHours(
                                                                                                        m.avg_close_time_hours,
                                                                                                )}
                                                                                        </td>
                                                                                </tr>
                                                                        ),
                                                                )}
                                                        </tbody>
                                                </table>
                                        </div>
                                </>
                        )}
                </div>
        );
}
