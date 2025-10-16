import axios from 'axios';
import { Post, CreatePostRequest, Comment, CreateCommentRequest } from '../types';

const API_URL = 'http://localhost:8080/api/posts';

export const postService = {
  // Obtener todos los posts
  async getAllPosts(): Promise<Post[]> {
    const response = await axios.get<Post[]>(API_URL);
    return response.data;
  },

  // Crear un nuevo post
  async createPost(data: CreatePostRequest, userId: number): Promise<Post> {
    const response = await axios.post<Post>(API_URL, data, {
      headers: {
        'X-User-ID': userId.toString()
      }
    });
    return response.data;
  },

  // Obtener un post por ID
  async getPostById(id: number): Promise<Post> {
    const response = await axios.get<Post>(`${API_URL}/${id}`);
    return response.data;
  },

  // Eliminar un post
  async deletePost(id: number, userId: number): Promise<void> {
    await axios.delete(`${API_URL}/${id}`, {
      headers: {
        'X-User-ID': userId.toString()
      }
    });
  },

  // Obtener comentarios de un post
  async getComments(postId: number): Promise<Comment[]> {
    const response = await axios.get<Comment[]>(`${API_URL}/${postId}/comments`);
    return response.data;
  },

  // Crear comentario
  async createComment(postId: number, data: CreateCommentRequest, userId: number): Promise<Comment> {
    const response = await axios.post<Comment>(
      `${API_URL}/${postId}/comments`,
      data,
      {
        headers: {
          'X-User-ID': userId.toString()
        }
      }
    );
    return response.data;
  }
};

// Eliminar comentario
export const deleteComment = async (postId: number, commentId: number, userId: number) => {
    return axios.delete(`${API_URL}/${postId}/comments/${commentId}`, {
        headers: { 'X-User-ID': userId.toString() }
    });
};