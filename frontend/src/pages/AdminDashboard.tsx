import React, { useEffect, useState } from 'react';
import api from '../api/axios';
import { ShieldCheck, CheckCircle } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export default function AdminDashboard() {
    const [transactions, setTransactions] = useState([]);
    const navigate = useNavigate();

    const fetchTxs = async () => {
        try {
            const { data } = await api.get('/admin/transactions');
            setTransactions(data);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchTxs();
    }, []);

    const handleComplete = async (txId: number) => {
        try {
            await api.post(`/admin/transactions/${txId}/complete`);
            fetchTxs(); // refresh
        } catch (err) {
            alert("Verification failed or not ready.");
        }
    };

    return (
        <div className="w-full animate-fade-in pt-10">
            <div className="flex items-center gap-3 mb-8">
                <ShieldCheck className="w-10 h-10 text-emerald-500" />
                <div>
                    <h1 className="text-4xl text-gradient font-bold">Admin Central</h1>
                    <p className="text-slate-400">Oversee escrow transactions and verify proofs.</p>
                </div>
            </div>

            <div className="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden shadow-2xl">
                <table className="w-full text-left text-sm text-slate-300">
                    <thead className="bg-slate-950 text-slate-400 border-b border-slate-800 uppercase text-xs font-semibold">
                        <tr>
                            <th className="px-6 py-4">TX ID</th>
                            <th className="px-6 py-4">Ticket</th>
                            <th className="px-6 py-4">Buyer / Seller</th>
                            <th className="px-6 py-4">Status</th>
                            <th className="px-6 py-4 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800/50">
                        {transactions.map((tx: any) => (
                            <tr key={tx.id} className="hover:bg-slate-800/20 transition-colors">
                                <td className="px-6 py-4 font-mono text-slate-500">#{tx.id}</td>
                                <td className="px-6 py-4 font-medium text-white">{tx.ticket?.title}</td>
                                <td className="px-6 py-4">
                                    <div className="flex flex-col gap-1 text-xs">
                                        <span className="bg-blue-500/10 text-blue-400 px-2 py-0.5 rounded w-max border border-blue-500/20">
                                            B: {tx.buyer?.username}
                                        </span>
                                        <span className="bg-orange-500/10 text-orange-400 px-2 py-0.5 rounded w-max border border-orange-500/20">
                                            S: {tx.seller?.username}
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`px-3 py-1 rounded-full text-xs font-medium border ${tx.status === 'Completed' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' :
                                            tx.status === 'Verifying' ? 'bg-amber-500/10 text-amber-400 border-amber-500/20' :
                                                'bg-slate-500/10 text-slate-400 border-slate-500/20'
                                        }`}>
                                        {tx.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right">
                                    <div className="flex justify-end gap-2">
                                        <button
                                            onClick={() => navigate(`/checkout/${tx.id}?chat=true`)}
                                            className="bg-slate-800 hover:bg-slate-700 text-white px-3 py-1.5 rounded-lg border border-slate-700 transition-colors"
                                        >
                                            Open Chat Room
                                        </button>
                                        {tx.status === 'Verifying' && (
                                            <button
                                                onClick={() => handleComplete(tx.id)}
                                                className="bg-emerald-600 hover:bg-emerald-500 text-white px-3 py-1.5 rounded-lg flex items-center gap-1 transition-colors"
                                            >
                                                <CheckCircle className="w-4 h-4" /> Finalize
                                            </button>
                                        )}
                                    </div>
                                </td>
                            </tr>
                        ))}
                        {transactions.length === 0 && (
                            <tr>
                                <td colSpan={5} className="px-6 py-10 text-center text-slate-500">
                                    No active escrow transactions.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
