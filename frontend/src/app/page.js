'use client';

import { useState, useEffect } from 'react';
import dynamic from 'next/dynamic';

const MapWithNoSSR = dynamic(() => import('../components/Map'), {
    ssr: false,
    loading: () => <p>Loading Map...</p>
});

export default function Home() {
    const [activeDistrictName, setActiveDistrictName] = useState(null); // Filter state
    const [incidents, setIncidents] = useState([]);
    const [stats, setStats] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
                const [resIncidents, resStats] = await Promise.all([
                    fetch(`${apiUrl}/incidents`),
                    fetch(`${apiUrl}/stats/district`)
                ]);

                const dataIncidents = await resIncidents.json();
                const dataStats = await resStats.json();

                setIncidents(dataIncidents || []);
                setStats(dataStats || []);
            } catch (error) {
                console.error("Failed to fetch data:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    // Filter Logic: Filter incidents if a district is selected
    const filteredIncidents = activeDistrictName
        ? incidents.filter(i => i.district?.name === activeDistrictName)
        : incidents;

    return (
        <div className="dashboard-container">
            {/* Sidebar: District Menu */}
            <aside className="sidebar">
                <div className="sidebar-header">
                    <h2>WASPADA</h2>
                    <p className="subtitle">Bandung Crime Monitor</p>
                </div>

                <div className="list-title">Districts</div>

                <div className="district-list">
                    <div
                        className={`district-item ${activeDistrictName === null ? 'active' : ''}`}
                        onClick={() => setActiveDistrictName(null)}
                    >
                        <span>All Districts</span>
                    </div>
                    {stats.map(stat => (
                        <div
                            key={stat.district_name}
                            className={`district-item ${activeDistrictName === stat.district_name ? 'active' : ''}`}
                            onClick={() => setActiveDistrictName(stat.district_name)}
                        >
                            <span>{stat.district_name}</span>
                            <span className="count-badge">
                                {stat.count}
                            </span>
                        </div>
                    ))}
                </div>
            </aside>

            {/* Main Content: Map */}
            <main className="main-content">
                <header className="header">
                    <div className="header-content">
                        <div style={{ fontWeight: '600', letterSpacing: '1px' }}>
                            {activeDistrictName ? `LIVE MAP: ${activeDistrictName.toUpperCase()}` : 'LIVE MAP: ALL BANDUNG'}
                        </div>
                        <div className="live-indicator">
                            <div className="pulsating-dot"></div>
                            <span className="live-text">{loading ? 'Connecting...' : 'Live System'}</span>
                        </div>
                    </div>
                </header>
                <div className="map-container">
                    <MapWithNoSSR
                        incidents={filteredIncidents}
                        stats={stats}
                        onDistrictClick={setActiveDistrictName}
                    />
                </div>
            </main>

            {/* Right Panel: Incident Feed (New Feature) */}
            <aside className="feed-panel">
                <div className="list-title">Latest Incidents</div>
                <div className="feed-list">
                    {filteredIncidents.length === 0 ? (
                        <div className="feed-item" style={{ textAlign: 'center', opacity: 0.6 }}>No incidents found.</div>
                    ) : (
                        filteredIncidents.map(incident => (
                            <div key={incident.id} className="feed-item">
                                <div className="feed-category">{incident.category || 'Uncategorized'}</div>
                                <a href={incident.source_url} target="_blank" rel="noreferrer" className="feed-title">
                                    {incident.title}
                                </a>
                                <p className="feed-desc">{incident.description ? incident.description.substring(0, 100) + '...' : 'No description available.'}</p>
                                <div className="feed-meta">
                                    <span className="feed-date">{new Date(incident.incident_date).toLocaleDateString()}</span>
                                    <span className="feed-district">{incident.district?.name}</span>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </aside>
        </div>
    );
}
