import React from 'react';
import { Container } from 'semantic-ui-react';

import Riddle from './components/Riddle';

const App = () => {
  return (
    <div>
        <Container>
            <Riddle/>
        </Container>
    </div>
  );
};

export default App;