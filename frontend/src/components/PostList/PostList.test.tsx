import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { PostList } from './PostList';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('PostList Component', () => {
  const mockPosts = [
    {
      id: 1,
      title: 'Mi primer post',
      content: 'Este es el contenido del primer post',
      user_id: 1,
      username: 'testuser',
      created_at: '2025-01-01T10:00:00Z'
    },
    {
      id: 2,
      title: 'Post de otro usuario',
      content: 'Este es de otro usuario',
      user_id: 2,
      username: 'otheruser',
      created_at: '2025-01-02T10:00:00Z'
    }
  ];

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renderiza la lista de posts correctamente', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockPosts });

    render(<PostList currentUserId={1} />);

    // Esperar a que cargue
    await waitFor(() => {
      expect(screen.getByText('Mi primer post')).toBeInTheDocument();
      expect(screen.getByText('Post de otro usuario')).toBeInTheDocument();
    });

    // Verificar que se muestre el contenido
    expect(screen.getByText('Este es el contenido del primer post')).toBeInTheDocument();
    expect(screen.getByText(/por @testuser/)).toBeInTheDocument();
    expect(screen.getByText(/por @otheruser/)).toBeInTheDocument();
  });

  test('muestra "No hay posts" cuando la lista está vacía', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: [] });

    render(<PostList currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText(/no hay posts/i)).toBeInTheDocument();
    });
  });

  test('muestra botón eliminar solo para posts propios', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockPosts });

    render(<PostList currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText('Mi primer post')).toBeInTheDocument();
    });

    // Debe haber un botón eliminar (para el post del usuario 1)
    const deleteButtons = screen.getAllByText('Eliminar');
    expect(deleteButtons).toHaveLength(1);

    // El post del usuario 2 no debe tener botón eliminar
  });

  test('elimina un post cuando se hace click en eliminar', async () => {
    mockedAxios.get.mockResolvedValueOnce({ data: mockPosts });
    mockedAxios.delete.mockResolvedValueOnce({ data: {} });
    mockedAxios.get.mockResolvedValueOnce({ data: [] }); // Segunda llamada después de eliminar

    window.confirm = jest.fn(() => true); // Mock de confirm

    render(<PostList currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText('Mi primer post')).toBeInTheDocument();
    });

    // Click en eliminar
    const deleteButton = screen.getByText('Eliminar');
    fireEvent.click(deleteButton);

    // Verificar que se llamó a delete con los parámetros correctos
    await waitFor(() => {
      expect(mockedAxios.delete).toHaveBeenCalledWith(
        'http://localhost:8080/api/posts/1',
        {
          headers: {
            'X-User-ID': '1'
          }
        }
      );
    });
  });

  test('muestra error cuando falla cargar posts', async () => {
    mockedAxios.get.mockRejectedValueOnce({
      response: {
        data: {
          error: 'Error en el servidor'
        }
      }
    });

    render(<PostList currentUserId={1} />);

    await waitFor(() => {
      expect(screen.getByText('Error al cargar posts')).toBeInTheDocument();
    });
  });
});