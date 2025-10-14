import React, { useState } from 'react';
import { postService } from '../../services/postService';
import './CreatePost.css';

interface CreatePostProps {
  userId: number;
  onPostCreated: () => void;
}

export const CreatePost: React.FC<CreatePostProps> = ({ userId, onPostCreated }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await postService.createPost({ title, content }, userId);
      
      // Limpiar formulario
      setTitle('');
      setContent('');
      
      // Notificar que se creó el post
      onPostCreated();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al crear post');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="create-post-container">
      <h2>Crear Nuevo Post</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Título:</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Escribe un título..."
            required
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label>Contenido:</label>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="¿Qué quieres compartir?"
            rows={5}
            required
            disabled={loading}
          />
        </div>

        {error && <div className="error-message">{error}</div>}

        <button type="submit" disabled={loading}>
          {loading ? 'Publicando...' : 'Publicar Post'}
        </button>
      </form>
    </div>
  );
};