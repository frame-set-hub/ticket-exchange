import React, { useState } from 'react';
import { useTickets, useCreateTicket } from '../features/ticket/hooks/useTickets';
import TicketCard from '../components/TicketCard';
import { Search, Filter, RefreshCw } from 'lucide-react';
import { useToast } from '../components/Toast';

export default function Home() {
    const { tickets, loading, fetchTickets } = useTickets();
    const { createTicket, loading: isSelling } = useCreateTicket();
    const { toast } = useToast();

    // Filters
    const [title, setTitle] = useState('');
    const [category, setCategory] = useState('');
    const [minPrice, setMinPrice] = useState('');
    const [maxPrice, setMaxPrice] = useState('');

    // Sell Ticket Modal
    const [isSellModalOpen, setIsSellModalOpen] = useState(false);
    const [sellForm, setSellForm] = useState({ title: '', venue: '', price: '', category: 'Concert', description: '' });

    const handleSell = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await createTicket({
                ...sellForm,
                price: parseFloat(sellForm.price),
            });
            setIsSellModalOpen(false);
            setSellForm({ title: '', venue: '', price: '', category: 'Concert', description: '' });
            fetchTickets();
            toast('success', 'Ticket Listed!', 'Your ticket is now live on the marketplace.');
        } catch {
            toast('error', 'Listing Failed', 'Could not create your ticket. Please try again.');
        }
    };

    const handleFilter = () => {
        fetchTickets({ title, category, min_price: minPrice, max_price: maxPrice });
    };

    return (
        <div className="w-full animate-fade-in pt-10">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-end mb-8 gap-4">
                <div>
                    <h1 className="text-4xl text-gradient font-bold mb-2">Marketplace</h1>
                    <p className="text-slate-400">Discover and securely buy secondary tickets.</p>
                </div>
                <button
                    onClick={() => setIsSellModalOpen(true)}
                    className="bg-emerald-600 hover:bg-emerald-500 text-white px-4 py-2 rounded-lg font-medium transition-colors shadow-lg border border-emerald-500 shadow-emerald-900/20"
                >
                    + Sell Ticket
                </button>
            </div>

            {/* Sell Ticket Modal */}
            {isSellModalOpen && (
                <div className="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
                    <div className="glass-panel w-full max-w-md p-6 rounded-2xl">
                        <h2 className="text-2xl font-bold text-white mb-4">Sell a Ticket</h2>
                        <form onSubmit={handleSell} className="space-y-4">
                            <div>
                                <label className="block text-sm text-slate-300 mb-1">Title</label>
                                <input required value={sellForm.title} onChange={e => setSellForm({ ...sellForm, title: e.target.value })} className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-white focus:border-blue-500 outline-none" placeholder="Coldplay BKK 2026" />
                            </div>
                            <div>
                                <label className="block text-sm text-slate-300 mb-1">Venue</label>
                                <input required value={sellForm.venue} onChange={e => setSellForm({ ...sellForm, venue: e.target.value })} className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-white focus:border-blue-500 outline-none" placeholder="Rajamangala Stadium" />
                            </div>
                            <div className="flex gap-4">
                                <div className="flex-1">
                                    <label className="block text-sm text-slate-300 mb-1">Category</label>
                                    <select value={sellForm.category} onChange={e => setSellForm({ ...sellForm, category: e.target.value })} className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-white focus:border-blue-500 outline-none appearance-none">
                                        <option value="Concert">Concert</option>
                                        <option value="Sports">Sports</option>
                                        <option value="Theater">Theater</option>
                                    </select>
                                </div>
                                <div className="flex-1">
                                    <label className="block text-sm text-slate-300 mb-1">Price ($)</label>
                                    <input required type="number" min="1" step="0.01" value={sellForm.price} onChange={e => setSellForm({ ...sellForm, price: e.target.value })} className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-white focus:border-blue-500 outline-none" placeholder="150" />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm text-slate-300 mb-1">Description</label>
                                <textarea value={sellForm.description} onChange={e => setSellForm({ ...sellForm, description: e.target.value })} className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-white focus:border-blue-500 outline-none" rows={3} placeholder="Seat A12, perfect view." />
                            </div>
                            <div className="flex justify-end gap-3 pt-4">
                                <button type="button" onClick={() => setIsSellModalOpen(false)} className="px-4 py-2 text-slate-300 hover:text-white transition-colors">Cancel</button>
                                <button type="submit" disabled={isSelling} className="bg-emerald-600 hover:bg-emerald-500 text-white px-6 py-2 rounded-lg font-medium transition-colors disabled:opacity-50">
                                    {isSelling ? 'Listing...' : 'Confirm'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            <div className="glass-panel p-4 rounded-xl mb-8 flex flex-col md:flex-row gap-4 items-end">
                <div className="flex-1 w-full relative">
                    <label className="text-xs font-medium text-slate-400 mb-1 block">Search Title</label>
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
                        <input
                            value={title} onChange={(e) => setTitle(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg pl-9 pr-4 py-2.5 text-white placeholder-slate-500 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                            placeholder="Searching for Concerts..."
                        />
                    </div>
                </div>

                <div className="w-full md:w-48">
                    <label className="text-xs font-medium text-slate-400 mb-1 block">Category</label>
                    <select
                        value={category} onChange={(e) => setCategory(e.target.value)}
                        className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-blue-500 appearance-none"
                    >
                        <option value="">All Categories</option>
                        <option value="Concert">Concert</option>
                        <option value="Sports">Sports</option>
                        <option value="Theater">Theater</option>
                    </select>
                </div>

                <div className="flex gap-2 w-full md:w-auto">
                    <div className="w-24">
                        <label className="text-xs font-medium text-slate-400 mb-1 block">Min $</label>
                        <input type="number" value={minPrice} onChange={(e) => setMinPrice(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2.5 text-white placeholder-slate-500 focus:outline-none focus:border-blue-500" placeholder="0" />
                    </div>
                    <div className="w-24">
                        <label className="text-xs font-medium text-slate-400 mb-1 block">Max $</label>
                        <input type="number" value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2.5 text-white placeholder-slate-500 focus:outline-none focus:border-blue-500" placeholder="Any" />
                    </div>
                </div>

                <button onClick={handleFilter} className="bg-slate-800 hover:bg-slate-700 border border-slate-600 text-slate-200 px-4 py-2.5 rounded-lg flex items-center justify-center gap-2 transition-colors md:w-auto w-full">
                    <Filter className="w-4 h-4" /> Filter
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {loading ? (
                    <div className="col-span-full flex justify-center py-20 text-slate-500 items-center gap-2">
                        <RefreshCw className="w-5 h-5 animate-spin" /> Loading Tickets...
                    </div>
                ) : tickets.length === 0 ? (
                    <div className="col-span-full text-center py-20 text-slate-500 bg-slate-900/50 rounded-2xl border border-slate-800 border-dashed">
                        No tickets found securely match your criteria.
                    </div>
                ) : (
                    tickets.map((ticket) => <TicketCard key={ticket.id} ticket={ticket} />)
                )}
            </div>
        </div>
    );
}
