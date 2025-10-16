import React, { useEffect, useState, useCallback } from 'react';
import { postService, deleteComment } from '../../services/postService';
import { Comment } from '../../types';
import './CommentList.css';

interface CommentListProps {
    postId: number;
    currentUserId: number;
    refreshTrigger?: number;
    onCommentDeleted?: (commentId: number) => void;
}

const CommentList: React.FC<CommentListProps> = ({
    postId,
    currentUserId,
    refreshTrigger,
    onCommentDeleted
}) => {
    const [comments, setComments] = useState<Comment[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');

    // Memoizamos loadComments para evitar warning de ESLint
    const loadComments = useCallback(async () => {
        try {
            setLoading(true);
            const data = await postService.getComments(postId);
            setComments(data);
            setError('');
        } catch (err: any) {
            setError('Error al cargar comentarios');
        } finally {
            setLoading(false);
        }
    }, [postId]);

    useEffect(() => {
        loadComments();
    }, [loadComments, refreshTrigger]);

    const handleDelete = async (commentId: number) => {
        try {
            await deleteComment(postId, commentId, currentUserId);

            // Actualizar lista de comentarios usando función de estado
            setComments(prev => prev.filter(c => c.id !== commentId));

            if (onCommentDeleted) onCommentDeleted(commentId);

            // Mostrar mensaje de éxito
            setSuccessMessage('Comentario eliminado exitosamente');
            setTimeout(() => setSuccessMessage(''), 3000);

        } catch (err: any) {
            console.error("Error deleting comment:", err.response?.data || err.message);
            alert("Error al eliminar el comentario");
        }
    };

    if (loading) return <div className="comments-loading">Cargando comentarios...</div>;
    if (error) return <div className="comments-error">{error}</div>;
    if (comments.length === 0) return <div className="no-comments">No hay comentarios todavía. ¡Sé el primero en comentar!</div>;

    return (
        <div className="comment-list">
            <h3>Comentarios ({comments.length})</h3>

            {successMessage && <div className="success-message">{successMessage}</div>}

            {comments.map(comment => (
                <div key={comment.id} className="comment-card">
                    <div className="comment-header">
                        <span className="comment-author">@{comment.username}</span>
                        <span className="comment-date">{new Date(comment.created_at).toLocaleDateString()}</span>
                        {comment.user_id === currentUserId && (
                            <button
                                className="comment-delete-btn"
                                onClick={() => handleDelete(comment.id)}
                            >
                                Eliminar
                            </button>
                        )}
                    </div>
                    <p className="comment-content">{comment.content}</p>
                </div>
            ))}
        </div>
    );
};

export default CommentList;