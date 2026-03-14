import { useState, useEffect, useRef, useCallback } from 'react';
import type { Message } from '../../../domains/chat/entities/Message';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

const WS_BASE = 'ws://localhost:8080/api/chat/ws';

export function useChat(token: string | null, transactionId: number | undefined) {
  const [messages, setMessages] = useState<Message[]>([]);
  const { chatRepository } = getServices();
  const wsRef = useRef<WebSocket | null>(null);

  // 1. Fetch history via REST on mount
  const fetchHistory = useCallback(async () => {
    if (!transactionId) return;
    try {
      const data = await chatRepository.getMessages(transactionId);
      setMessages(data ?? []);
    } catch {
      // no messages yet
    }
  }, [transactionId]);

  useEffect(() => {
    fetchHistory();
  }, [fetchHistory]);

  // 2. Connect WebSocket for real-time
  useEffect(() => {
    if (!transactionId || !token) return;

    const url = `${WS_BASE}/${transactionId}?token=${token}`;
    const ws = new WebSocket(url);

    ws.onopen = () => {
      console.log(`[WS] Connected to room tx:${transactionId}`);
    };

    ws.onmessage = (event) => {
      try {
        const msg: Message = JSON.parse(event.data);
        setMessages((prev) => {
          // Deduplicate by id
          if (prev.some((m) => m.id === msg.id)) return prev;
          return [...prev, msg];
        });
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = () => {
      console.log(`[WS] Disconnected from room tx:${transactionId}`);
    };

    wsRef.current = ws;

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [transactionId, token]);

  // 3. Send message via WebSocket
  const sendMessage = (content: string) => {
    if (!content) return;
    const ws = wsRef.current;
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ content }));
    }
  };

  return { messages, sendMessage };
}
