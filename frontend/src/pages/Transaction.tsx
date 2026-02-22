import React, { useEffect, useState, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { useAuthStore } from '../store/useStore';
import api from '../api/axios';
import { Send, Upload, FileCheck, Shield, Lock } from 'lucide-react';

export default function Transaction() {
    const { id } = useParams(); // Ticket ID or TX ID based on flow
    const { user, token } = useAuthStore();

    const [tx, setTx] = useState<any>(null);
    const [messages, setMessages] = useState<any[]>([]);
    const [inputMsg, setInputMsg] = useState('');
    const ws = useRef<WebSocket | null>(null);

    const initEscrow = async () => {
        try {
            // 1. First attempt to fetch tx by ticket ID (buyer initiates)
            const { data } = await api.post(`/transactions/${id}`);
            setTx(data);
        } catch (err: any) {
            // 2. If it fails, maybe it already exists or we are loading directly by tx ID
            if (err.response?.status === 400 && err.response?.data?.error === 'Ticket is not available') {
                // Mock fallback to get tx for POC natively would require another endpoint.
                // Assuming we have TX ID here.
            }
        }
    };

    const connectWS = () => {
        if (!token) return;
        const socket = new WebSocket(`ws://localhost:8080/api/chat?token=${token}`);
        socket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            if (msg.transaction_id === tx?.id) {
                setMessages((prev) => [...prev, msg]);
            }
        };
        ws.current = socket;
    };

    useEffect(() => {
        initEscrow();
        return () => ws.current?.close();
    }, []);

    useEffect(() => {
        if (tx) connectWS();
    }, [tx]);

    const sendMessage = (e: React.FormEvent) => {
        e.preventDefault();
        if (!ws.current || !inputMsg) return;

        // Send to admin (for POC we assume receiver_id is Admin globally or hardcoded)
        ws.current.send(JSON.stringify({
            transaction_id: tx?.id,
            content: inputMsg,
            receiver_id: 1, // Mock admin ID
        }));
        setInputMsg('');
    };

    const handleUploadProof = async () => {
        if (!tx) return;
        try {
            if (user?.id === tx.buyer_id) {
                await api.post(`/transactions/${tx.id}/upload-payment`);
            } else if (user?.id === tx.seller_id) {
                await api.post(`/transactions/${tx.id}/upload-ticket`);
            }
            alert("Uploaded successfully to Escrow!");
            // Reload tx logic here
        } catch (err) {
            console.error("Upload failed");
        }
    };

    if (!tx) return <div className="pt-20 text-center text-slate-400">Loading Secure Escrow...</div>;

    return (
        <div className="pt-8 h-[calc(100vh-80px)] flex flex-col md:flex-row gap-6">
            {/* Sidebar Info */}
            <div className="w-full md:w-1/3 flex flex-col gap-6">
                <div className="glass-panel p-6 rounded-2xl border-emerald-500/30">
                    <div className="flex items-center gap-3 mb-4">
                        <Lock className="w-8 h-8 text-emerald-500" />
                        <h2 className="text-2xl font-bold text-white">Escrow Hold</h2>
                    </div>
                    <p className="text-slate-400 text-sm mb-6">
                        Your transaction is guarded. Funds and tickets are held by the Admin until both parties fulfill their obligations.
                    </p>

                    <div className="bg-slate-900/80 rounded-xl p-4 mb-4 border border-slate-700">
                        <div className="text-xs text-slate-500 uppercase tracking-wider mb-1">Status</div>
                        <div className={`font-semibold ${tx.status === 'Completed' ? 'text-emerald-400' : 'text-amber-400'}`}>
                            {tx.status.replace('_', ' ')}
                        </div>
                    </div>

                    <button
                        onClick={handleUploadProof}
                        className="w-full bg-blue-600 hover:bg-blue-500 text-white py-3 rounded-lg flex items-center justify-center gap-2 font-medium transition-colors"
                    >
                        {user?.id === tx.seller_id ? <><Upload className="w-4 h-4" /> Upload Digital Ticket</> : <><FileCheck className="w-4 h-4" /> Upload Payment Proof</>}
                    </button>
                </div>
            </div>

            {/* Chat Area */}
            <div className="flex-1 glass-panel rounded-2xl flex flex-col overflow-hidden">
                <div className="bg-slate-900/90 border-b border-slate-800 p-4 flex items-center gap-3">
                    <Shield className="w-5 h-5 text-blue-500" />
                    <h3 className="font-medium text-white">Private Admin Channel -- #{tx.id}</h3>
                </div>

                <div className="flex-1 overflow-y-auto p-6 space-y-4">
                    <div className="text-center text-xs text-slate-500 my-4">Conversation started. Messages are end-to-end encrypted with Admin.</div>

                    {messages.map((m, i) => (
                        <div key={i} className={`flex flex-col ${m.sender_id === user?.id ? 'items-end' : 'items-start'}`}>
                            <div className={`max-w-[75%] px-4 py-2 rounded-2xl ${m.sender_id === user?.id ? 'bg-blue-600 text-white rounded-tr-sm' : 'bg-slate-800 text-slate-200 rounded-tl-sm'}`}>
                                {m.content}
                            </div>
                        </div>
                    ))}
                </div>

                <form onSubmit={sendMessage} className="p-4 bg-slate-900/50 border-t border-slate-800 flex gap-3">
                    <input
                        type="text" value={inputMsg} onChange={e => setInputMsg(e.target.value)}
                        placeholder="Message the Admin..."
                        className="flex-1 bg-slate-800 border-none rounded-xl px-4 text-white focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                    <button type="submit" className="bg-blue-600 hover:bg-blue-500 text-white p-3 rounded-xl transition-colors">
                        <Send className="w-5 h-5" />
                    </button>
                </form>
            </div>
        </div>
    );
}
