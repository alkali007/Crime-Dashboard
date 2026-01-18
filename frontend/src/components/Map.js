'use client';

import { MapContainer, TileLayer, Marker, Popup, CircleMarker } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.css';
import 'leaflet-defaulticon-compatibility';

const MapComponent = ({ incidents, stats }) => {
    const center = [-6.9175, 107.6191]; // Bandung Center

    return (
        <MapContainer center={center} zoom={13} style={{ height: '100%', width: '100%' }}>
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />

            {incidents.map((incident) => (
                incident.District && (
                    <Marker
                        key={incident.id}
                        position={[incident.District.latitude, incident.District.longitude]}
                    >
                        <Popup>
                            <strong>{incident.title}</strong><br />
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
                    color="red"
                    fillColor="#f03"
                    fillOpacity={0.3}
                >
                    <Popup>
                        {stat.district_name}: {stat.count} incidents
                    </Popup>
                </CircleMarker>
            ))}

        </MapContainer>
    );
};

export default MapComponent;
