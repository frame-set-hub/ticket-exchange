import { MapPin } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import type { Ticket } from '../domains/ticket/entities/Ticket';

export default function TicketCard({ ticket }: { ticket: Ticket }) {
    const navigate = useNavigate();

    return (
        <div className="glass-panel p-5 rounded-2xl flex flex-col hover:-translate-y-1 transition-transform cursor-pointer relative overflow-hidden group">
            <div className="absolute top-0 right-0 p-3">
                <span className="bg-blue-500/10 text-blue-400 text-xs font-semibold px-2 py-1 rounded-full border border-blue-500/20">
                    {ticket.category}
                </span>
            </div>

            <div className="mt-2 text-2xl font-bold text-white mb-1 group-hover:text-blue-400 transition-colors">
                {ticket.title}
            </div>

            <div className="flex items-center text-slate-400 text-sm mb-4">
                <MapPin className="w-4 h-4 mr-1 text-slate-500" />
                {ticket.venue}
            </div>

            <div className="flex-1" />

            <div className="flex items-center justify-between mt-4 p-4 rounded-xl bg-slate-900/50 border border-slate-700/50">
                <div>
                    <div className="text-sm text-slate-400 mb-0.5">Price</div>
                    <div className="text-xl font-bold text-emerald-400">${ticket.price.toFixed(2)}</div>
                </div>
                <button
                    onClick={() => navigate(`/checkout/${ticket.id}`)}
                    className="bg-blue-600 hover:bg-blue-500 text-white px-5 py-2.5 rounded-lg font-medium transition-colors border border-blue-500"
                >
                    Buy Now
                </button>
            </div>
            <div className="mt-4 flex items-center justify-between text-xs text-slate-500">
                <div>Seller: <span className="text-slate-300">{ticket.seller?.username || 'Unknown'}</span></div>
            </div>
        </div>
    );
}
