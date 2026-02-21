'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listsAPI } from '@/lib/api';
import type { UserList } from '@/lib/api';
import { Plus, Trash2, Edit2, List, Lock, Unlock, X } from 'lucide-react';
import toast from 'react-hot-toast';

export default function ListsPage() {
  const queryClient = useQueryClient();
  const [isCreating, setIsCreating] = useState(false);
  const [editingList, setEditingList] = useState<UserList | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    is_public: false,
  });

  const { data: lists, isLoading } = useQuery<UserList[]>({
    queryKey: ['user-lists'],
    queryFn: listsAPI.getLists,
  });

  const createMutation = useMutation({
    mutationFn: listsAPI.createList,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-lists'] });
      toast.success('List created successfully!');
      setIsCreating(false);
      setFormData({ name: '', description: '', is_public: false });
    },
    onError: () => toast.error('Failed to create list'),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: Parameters<typeof listsAPI.updateList>[1] }) =>
      listsAPI.updateList(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-lists'] });
      toast.success('List updated successfully!');
      setEditingList(null);
      setFormData({ name: '', description: '', is_public: false });
    },
    onError: () => toast.error('Failed to update list'),
  });

  const deleteMutation = useMutation({
    mutationFn: listsAPI.deleteList,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-lists'] });
      toast.success('List deleted successfully!');
    },
    onError: () => toast.error('Failed to delete list'),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      toast.error('List name is required');
      return;
    }

    if (editingList) {
      updateMutation.mutate({
        id: editingList.list_id,
        data: {
          name: formData.name,
          description: formData.description || undefined,
          is_public: formData.is_public,
        },
      });
    } else {
      createMutation.mutate({
        name: formData.name,
        description: formData.description || undefined,
        is_public: formData.is_public,
      });
    }
  };

  const handleEdit = (list: UserList) => {
    setEditingList(list);
    setFormData({
      name: list.name,
      description: list.description || '',
      is_public: list.is_public,
    });
    setIsCreating(true);
  };

  const handleCancel = () => {
    setIsCreating(false);
    setEditingList(null);
    setFormData({ name: '', description: '', is_public: false });
  };

  const handleDelete = (id: number, name: string) => {
    if (confirm(`Are you sure you want to delete "${name}"?`)) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading lists...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">My Lists</h1>
          <p className="text-gray-600 mt-2">Organize problems into custom lists</p>
        </div>
        {!isCreating && (
          <button onClick={() => setIsCreating(true)} className="btn-primary flex items-center gap-2">
            <Plus className="w-5 h-5" />
            New List
          </button>
        )}
      </div>

      {isCreating && (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900">
              {editingList ? 'Edit List' : 'Create New List'}
            </h2>
            <button onClick={handleCancel} className="text-gray-500 hover:text-gray-700">
              <X className="w-5 h-5" />
            </button>
          </div>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                List Name *
              </label>
              <input
                type="text"
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="input"
                placeholder="e.g., Dynamic Programming Basics"
                required
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                id="description"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="input"
                rows={3}
                placeholder="Optional description for your list"
              />
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="is_public"
                checked={formData.is_public}
                onChange={(e) => setFormData({ ...formData, is_public: e.target.checked })}
                className="w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
              />
              <label htmlFor="is_public" className="text-sm font-medium text-gray-700">
                Make this list public
              </label>
            </div>

            <div className="flex gap-3">
              <button
                type="submit"
                className="btn-primary"
                disabled={createMutation.isPending || updateMutation.isPending}
              >
                {editingList ? 'Update List' : 'Create List'}
              </button>
              <button type="button" onClick={handleCancel} className="btn-secondary">
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {lists && lists.length > 0 ? (
          lists.map((list) => (
            <div key={list.list_id} className="card hover:shadow-lg transition-shadow">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <List className="w-5 h-5 text-primary-600" />
                  <h3 className="font-bold text-lg text-gray-900">{list.name}</h3>
                </div>
                <div className="flex items-center gap-1">
                  {list.is_public ? (
                    <span title="Public">
                      <Unlock className="w-4 h-4 text-green-600" />
                    </span>
                  ) : (
                    <span title="Private">
                      <Lock className="w-4 h-4 text-gray-400" />
                    </span>
                  )}
                </div>
              </div>

              {list.description && (
                <p className="text-sm text-gray-600 mb-4 line-clamp-2">{list.description}</p>
              )}

              <div className="flex items-center justify-between mb-4">
                <div className="text-sm text-gray-700">
                  <span className="font-semibold">{list.total_items}</span> problems
                </div>
                <div className="text-sm text-gray-700">
                  <span className="font-semibold">{list.completed}</span> completed
                </div>
              </div>

              {list.total_items > 0 && (
                <div className="mb-4">
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-primary-600 h-2 rounded-full"
                      style={{ width: `${(list.completed / list.total_items) * 100}%` }}
                    />
                  </div>
                </div>
              )}

              <div className="flex gap-2 pt-4 border-t border-gray-200">
                <button
                  onClick={() => handleEdit(list)}
                  className="flex-1 btn-secondary text-sm flex items-center justify-center gap-1"
                >
                  <Edit2 className="w-4 h-4" />
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(list.list_id, list.name)}
                  className="flex-1 btn-secondary text-sm flex items-center justify-center gap-1 text-red-600 hover:bg-red-50"
                  disabled={deleteMutation.isPending}
                >
                  <Trash2 className="w-4 h-4" />
                  Delete
                </button>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full text-center py-12">
            <List className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No lists yet</h3>
            <p className="text-gray-600 mb-4">
              Create your first list to organize problems by topic or difficulty
            </p>
            <button
              onClick={() => setIsCreating(true)}
              className="btn-primary inline-flex items-center gap-2"
            >
              <Plus className="w-5 h-5" />
              Create Your First List
            </button>
          </div>
        )}
      </div>
    </div>
  );
}