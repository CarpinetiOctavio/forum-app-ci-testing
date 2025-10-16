import React, { useEffect, useState } from 'react';
import { postService } from '../../services/postService';
import { Post } from '../../types';
import  CommentList  from '../CommentList/CommentList';
import { CommentForm } from '../CommentForm/CommentForm';
import './PostDetail.css';


interface PostDetailProps {
  postId: number;
  userId: number;
  onBack: () => void;
}

export const PostDetail: React.FC<PostDetailProps> = ({ postId, userId, onBack }) => {
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [refreshComments, setRefreshComments] = useState(0);

  useEffect(() => {
    const loadPost = async () => {
      try {
        setLoading(true);
        const data = await postService.getPostById(postId);
        setPost(data);
        setError('');
      } catch (err: any) {
        setError('Error al cargar el post');
      } finally {
        setLoading(false);
      }
    };

    loadPost();
  }, [postId]);

  const handleCommentCreated = () => {
    setRefreshComments(prev => prev + 1);
  };

  if (loading) {
    return <div className="post-detail-loading">Cargando post...</div>;
  }

  if (error || !post) {
    return (
      <div className="post-detail-error">
        <p>{error || 'Post no encontrado'}</p>
        <button onClick={onBack}>Volver</button>
      </div>
    );
  }

  return (
    <div className="post-detail-container">
      <button className="back-btn" onClick={onBack}>
        ← Volver
      </button>

      <div className="post-detail-card">
        <h1>{post.title}</h1>
        <div className="post-detail-meta">
          <span className="post-detail-author">Por @{post.username}</span>
          <span className="post-detail-date">
            {new Date(post.created_at).toLocaleDateString()}
          </span>
        </div>
        <p className="post-detail-content">{post.content}</p>
      </div>

      <CommentForm 
        postId={postId} 
        userId={userId} 
        onCommentCreated={handleCommentCreated} 
      />

      <CommentList
        postId={postId}
        currentUserId={userId}
        refreshTrigger={refreshComments}
      />
    </div>
  );
};