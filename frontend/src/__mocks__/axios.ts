// Mock completo de axios para tests
const axiosMock = {
    get: jest.fn(() => Promise.resolve({ data: {} })),
    post: jest.fn(() => Promise.resolve({ data: {} })),
    put: jest.fn(() => Promise.resolve({ data: {} })),
    delete: jest.fn(() => Promise.resolve({ data: {} })),
    create: jest.fn(),
  };
  
  // Configurar create para que devuelva el mismo mock
  axiosMock.create.mockReturnValue(axiosMock);
  
  export default axiosMock;