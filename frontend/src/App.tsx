import React, { useState } from 'react';
import { Login } from './components/Login/Login';
import { PostList } from './components/PostList/PostList';
import { CreatePost } from './components/CreatePost/CreatePost';
import { User } from './types';
import './App.css';

function App() {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [refreshPosts, setRefreshPosts] = useState(false);

  const handleLoginSuccess = (user: User) => {
    setCurrentUser(user);
  };

  const handleLogout = () => {
    setCurrentUser(null);
  };

  const handlePostCreated = () => {
    // Forzar refresh de posts
    setRefreshPosts(!refreshPosts);
  };

  // Si no está logueado, mostrar login
  if (!currentUser) {
    return <Login onLoginSuccess={handleLoginSuccess} />;
  }

  // Si está logueado, mostrar la app
  return (
    <div className="App">
      <header className="app-header">
        <h1>🚀 Mini Red Social</h1>
        <div className="user-info">
          <span>Hola, @{currentUser.username}</span>
          <button onClick={handleLogout} className="logout-btn">
            Cerrar Sesión
          </button>
        </div>
      </header>

      <main>
        <CreatePost userId={currentUser.id} onPostCreated={handlePostCreated} />
        <PostList currentUserId={currentUser.id} onRefresh={refreshPosts} />
      </main>
    </div>
  );
}

export default App;