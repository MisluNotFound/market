import { Outlet } from 'react-router-dom';
import styled from 'styled-components';

const AppContainer = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
`;

const MainContent = styled.main`
  flex: 1;
  padding: 20px 0;
`;

function App() {
  return (
    <AppContainer>
      <MainContent>
        <Outlet />
      </MainContent>
    </AppContainer>
  );
}

export default App;
