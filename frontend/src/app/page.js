'use client';

import { useState, useEffect } from 'react';
import dynamic from 'next/dynamic';

const MapWithNoSSR = dynamic(() => import('../components/Map'), {
    ssr: false,
    loading: () => <p>Loading Map...</p>
});

export default function Home() {
    const [activeDistrict, setActiveDistrict] = useState(null);
    const [incidents, setIncidents] = useState([]);
    const [stats, setStats] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Fetch Data from Go Backend
        // Assuming backend is running on 8080

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

    const filteredIncidents = activeDistrict
        ? incidents.filter(i => i.district_id === activeDistrict.id) // This assumes we have district ID in sidebar, might need adjustment
        : incidents;

    return (
        <div className="dashboard-container">
            <aside className="sidebar">
                <div className="sidebar-header">
                    <h2>WASPADA</h2>
                </div>

                <div className="list-title">Districts Overview</div>

                <div className="district-list">
                    <div
                        className={`district-item ${activeDistrict === null ? 'active' : ''}`}
                        onClick={() => setActiveDistrict(null)}
                    >
                        <span>All Districts</span>
                    </div>
                    {stats.map(stat => (
                        <div
                            key={stat.district_name}
                            className={`district-item ${activeDistrict && activeDistrict.id === stat.district_name ? 'active' : ''}`} // Logic tweak: stat name match
                            onClick={() => setActiveDistrict(stat.district_name)} // Note: Logic remains same, simple state
                        >
                            <span>{stat.district_name}</span>
                            <span className="count-badge">
                                {stat.count}
                            </span>
                        </div>
                    ))}
                </div>
            </aside>

            <main className="main-content">
                <header className="header">
                    <div className="header-content">
                        <div style={{ fontWeight: '600', letterSpacing: '1px' }}>LIVE INCIDENTS MAP</div>
                        <div className="live-indicator">
                            <div className="pulsating-dot"></div>
                            <span className="live-text">{loading ? 'Connecting...' : 'Live System'}</span>
                        </div>
                    </div>
                </header>
                <div className="map-container">
                    <MapWithNoSSR incidents={incidents} stats={stats} />
                </div>
            </main>
        </div>
    );
}
