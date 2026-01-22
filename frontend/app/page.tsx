"use client";

import { useState, useEffect } from "react";
import { Plus, Trash2, CheckCircle2, Circle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";

interface Todo {
  id: number;
  title: string;
  completed: boolean;
  created_at: string;
}

export default function TodoApp() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState("");
  const [loading, setLoading] = useState(true);

  // Use the window location to determine API URL or fallback to localhost
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:30002";

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    try {
      const res = await fetch(`${API_URL}/todos`);
      const data = await res.json();
      setTodos(data || []);
    } catch (error) {
      console.error("Failed to fetch todos:", error);
    } finally {
      setLoading(false);
    }
  };

  const addTodo = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTodo.trim()) return;

    try {
      const res = await fetch(`${API_URL}/todos`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: newTodo, completed: false }),
      });
      const data = await res.json();
      setTodos([data, ...todos]);
      setNewTodo("");
    } catch (error) {
      console.error("Failed to add todo:", error);
    }
  };

  const toggleTodo = async (todo: Todo) => {
    try {
      const res = await fetch(`${API_URL}/todos/${todo.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...todo, completed: !todo.completed }),
      });
      const updatedTodo = await res.json();
      setTodos(todos.map((t) => (t.id === todo.id ? updatedTodo : t)));
    } catch (error) {
      console.error("Failed to toggle todo:", error);
    }
  };

  const deleteTodo = async (id: number) => {
    try {
      await fetch(`${API_URL}/todos/${id}`, { method: "DELETE" });
      setTodos(todos.filter((t) => t.id !== id));
    } catch (error) {
      console.error("Failed to delete todo:", error);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8">
      <Card className="mx-auto max-w-md shadow-lg border-slate-200">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center flex items-center justify-center gap-2">
            <CheckCircle2 className="w-6 h-6 text-primary" />
            Task Master
          </CardTitle>
          <p className="text-sm text-muted-foreground text-center">
            Organize your day, one task at a time.
          </p>
        </CardHeader>
        <CardContent className="space-y-6">
          <form onSubmit={addTodo} className="flex gap-2">
            <Input
              placeholder="What needs to be done?"
              value={newTodo}
              onChange={(e) => setNewTodo(e.target.value)}
              className="flex-1"
            />
            <Button type="submit" size="icon">
              <Plus className="h-4 w-4" />
            </Button>
          </form>

          <ScrollArea className="h-[400px] pr-4">
            {loading ? (
              <p className="text-center text-slate-500 py-10">Loading tasks...</p>
            ) : todos.length === 0 ? (
              <div className="text-center py-10 space-y-2">
                <Circle className="w-8 h-8 mx-auto text-slate-300" />
                <p className="text-slate-500">No tasks yet. Add one above!</p>
              </div>
            ) : (
              <div className="space-y-3">
                {todos.map((todo) => (
                  <div
                    key={todo.id}
                    className={cn(
                      "flex items-center justify-between p-3 rounded-lg border transition-all duration-200",
                      todo.completed 
                        ? "bg-slate-50 border-slate-100 opacity-70" 
                        : "bg-white border-slate-200 shadow-sm hover:border-primary/30"
                    )}
                  >
                    <div className="flex items-center gap-3 flex-1">
                      <Checkbox
                        checked={todo.completed}
                        onCheckedChange={() => toggleTodo(todo)}
                      />
                      <span className={cn(
                        "text-sm font-medium transition-all",
                        todo.completed && "line-through text-slate-500"
                      )}>
                        {todo.title}
                      </span>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => deleteTodo(todo.id)}
                      className="text-slate-400 hover:text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </ScrollArea>
        </CardContent>
      </Card>
    </div>
  );
}
