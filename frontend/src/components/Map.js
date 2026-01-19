'use client';

import { MapContainer, TileLayer, Marker, Popup, CircleMarker } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.css';
import 'leaflet-defaulticon-compatibility';

const MapComponent = ({ incidents, stats, onDistrictClick }) => {
    const center = [-6.9175, 107.6191]; // Bandung Center

    return (
        <MapContainer center={center} zoom={13} style={{ height: '100%', width: '100%' }}>
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />

            {incidents.map((incident) => (
                incident.district && (
                    <Marker
                        key={incident.id}
                        position={[incident.district.latitude, incident.district.longitude]}
                    >
                        <Popup>
                            <strong>{incident.title}</strong><br />
                            <span style={{ fontSize: '12px', color: '#666' }}>{incident.category}</span><br />
                            {new Date(incident.incident_date).toLocaleDateString()}<br />
                            <a href={incident.source_url} target="_blank" rel="noreferrer">Read More</a>
                        </Popup>
                    </Marker>
                )
            ))}

            {/* Heatmap-like circles for stats */}
            {stats.map((stat) => (
                <CircleMarker
                    key={stat.district_name}
                    center={[stat.latitude, stat.longitude]}
                    radius={Math.sqrt(stat.count) * 10 + 5}
                    color="#FF4B4B"
                    fillColor="#FF4B4B"
                    fillOpacity={0.4}
                    eventHandlers={{
                        click: () => {
                            if (onDistrictClick) onDistrictClick(stat.district_name);
                        },
                    }}
                >
                    <Popup>
                        <strong style={{ fontSize: '14px' }}>{stat.district_name}</strong><br />
                        {stat.count} incidents<br />
                        <button
                            style={{ marginTop: '5px', cursor: 'pointer', background: '#eee', border: 'none', padding: '4px 8px', borderRadius: '4px' }}
                            onClick={() => onDistrictClick && onDistrictClick(stat.district_name)}
                        >
                            Filter This District
                        </button>
                    </Popup>
                </CircleMarker>
            ))}

        </MapContainer>
    );
};

export default MapComponent;
