import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import  CommentList from './CommentList';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('CommentList Component', () => {
  const mockComments = [
    {
      id: 1,
      post_id: 1,
      user_id: 1,
      username: 'testuser',
      content: 'Mi comentario',
      created_at: '2025-01-01T10:00:00Z'
    },
    {
      id: 2,
      post_id: 1,
      user_id: 2,
      username: 'otheruser',
      content: 'Otro comentario',
      created_at: '2025-01-02T10:00:00Z'
    }
  ];

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renderiza la lista de comentarios correctamente', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockComments });

    render(<CommentList postId={1} currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText('Mi comentario')).toBeInTheDocument();
      expect(screen.getByText('Otro comentario')).toBeInTheDocument();
    });

    expect(screen.getByText('Comentarios (2)')).toBeInTheDocument();
  });

  test('muestra "No hay comentarios" cuando está vacía', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: [] });

    render(<CommentList postId={1} currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText(/no hay comentarios todavía/i)).toBeInTheDocument();
    });
  });

  test('muestra botón eliminar solo para comentarios propios', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockComments });

    render(<CommentList postId={1} currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText('Mi comentario')).toBeInTheDocument();
    });

    // Solo debe haber 1 botón eliminar (para el comentario del usuario 1)
    const deleteButtons = screen.queryAllByText(/eliminar/i);
    expect(deleteButtons).toHaveLength(1);
  });

  test('elimina un comentario cuando se hace click en eliminar', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockComments });
    mockedAxios.delete.mockResolvedValueOnce({ data: {} });

    const mockOnCommentDeleted = jest.fn();

    render(
      <CommentList 
        postId={1} 
        currentUserId={1}
        onCommentDeleted={mockOnCommentDeleted}
      />
    );

    await waitFor(() => {
      expect(screen.getByText('Mi comentario')).toBeInTheDocument();
    });

    // Click en eliminar
    const deleteButton = screen.getByText(/eliminar/i);
    fireEvent.click(deleteButton);

    // Verificar que se llamó a delete
    await waitFor(() => {
      expect(mockedAxios.delete).toHaveBeenCalledWith(
        'http://localhost:8080/api/posts/1/comments/1',
        {
          headers: {
            'X-User-ID': '1'
          }
        }
      );
    });

    // Verificar que se llamó el callback
    expect(mockOnCommentDeleted).toHaveBeenCalledWith(1);
  });

  test('muestra error cuando falla cargar comentarios', async () => {
    mockedAxios.get.mockRejectedValueOnce({
      response: {
        data: {
          error: 'Error del servidor'
        }
      }
    });

    render(<CommentList postId={1} currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText('Error al cargar comentarios')).toBeInTheDocument();
    });
  });
});