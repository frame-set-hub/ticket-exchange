import React, { useEffect, useState } from 'react';
import api from '../api/axios';
import { Ticket, Trash2, RefreshCw } from 'lucide-react';

export default function MyTickets() {
    const [tickets, setTickets] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchMyTickets = async () => {
        setLoading(true);
        try {
            const { data } = await api.get('/tickets/my');
            setTickets(data);
        } catch (err) {
            console.error(err);
        }
        setLoading(false);
    };

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure you want to delete this ticket?")) return;
        try {
            await api.delete(`/tickets/${id}`);
            fetchMyTickets(); // Refresh
        } catch (err: any) {
            alert(err.response?.data?.error || "Failed to delete ticket");
        }
    };

    useEffect(() => {
        fetchMyTickets();
    }, []);

    return (
        <div className="w-full animate-fade-in pt-10">
            <div className="flex justify-between items-end mb-8">
                <div>
                    <h1 className="text-4xl text-gradient font-bold mb-2">My Tickets</h1>
                    <p className="text-slate-400">Manage all tickets you have listed on the platform.</p>
                </div>
            </div>

            <div className="glass-panel overflow-hidden rounded-2xl">
                <table className="w-full text-left border-collapse">
                    <thead>
                        <tr className="bg-slate-900/80 border-b border-slate-700/50">
                            <th className="p-4 text-slate-300 font-medium">Ticket ID</th>
                            <th className="p-4 text-slate-300 font-medium">Title</th>
                            <th className="p-4 text-slate-300 font-medium">Venue</th>
                            <th className="p-4 text-slate-300 font-medium">Price</th>
                            <th className="p-4 text-slate-300 font-medium">Status</th>
                            <th className="p-4 text-slate-300 font-medium text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr>
                                <td colSpan={6} className="text-center p-8 text-slate-500">
                                    <div className="flex justify-center items-center gap-2">
                                        <RefreshCw className="w-5 h-5 animate-spin" /> Fetching Tickets...
                                    </div>
                                </td>
                            </tr>
                        ) : tickets.length === 0 ? (
                            <tr>
                                <td colSpan={6} className="text-center p-12 text-slate-500 bg-slate-900/20">
                                    <Ticket className="w-12 h-12 mx-auto mb-3 opacity-20" />
                                    <p>You haven't listed any tickets yet.</p>
                                </td>
                            </tr>
                        ) : (
                            tickets.map((t: any) => (
                                <tr key={t.id} className="border-b border-slate-800 hover:bg-slate-800/50 transition-colors">
                                    <td className="p-4 font-mono text-slate-400 text-sm">#{t.id}</td>
                                    <td className="p-4 text-white font-medium">{t.title}</td>
                                    <td className="p-4 text-slate-400">{t.venue}</td>
                                    <td className="p-4 text-emerald-400 font-semibold">${t.price}</td>
                                    <td className="p-4">
                                        <span className={`px-2.5 py-1 text-xs font-medium rounded-full ${t.status === 'Available' ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20' :
                                                t.status === 'Pending' ? 'bg-orange-500/10 text-orange-400 border border-orange-500/20' :
                                                    t.status === 'Waiting_Payment' ? 'bg-yellow-500/10 text-yellow-400 border border-yellow-500/20' :
                                                        t.status === 'Verifying' ? 'bg-purple-500/10 text-purple-400 border border-purple-500/20' :
                                                            'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                                            }`}>
                                            {t.status.replace('_', ' ')}
                                        </span>
                                    </td>
                                    <td className="p-4 text-right">
                                        {t.status === 'Available' ? (
                                            <button
                                                onClick={() => handleDelete(t.id)}
                                                className="p-2 text-rose-400 hover:bg-rose-500/10 rounded-lg transition-colors"
                                                title="Delete Ticket"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        ) : (
                                            <span className="text-xs text-slate-500 italic">Locked in Escrow</span>
                                        )}
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
