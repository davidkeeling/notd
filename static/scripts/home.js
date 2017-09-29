import ReactDOM from 'react-dom';

import Homepage from './scene';

ReactDOM.render(
  <Homepage population={population} isAdmin={isAdmin} />,
  document.getElementById("mount")
);
