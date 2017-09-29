import React from 'react';

function getRandomInt(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min)) + min;
}

export default class Homepage extends React.Component {
  constructor (props) {
    super(props);

    this.getInitialPosition = () => props.population.map(({ height, width, top, left, name }) => {
      return {
        height: height * .5,
        width: width * .5,
        top: top * .5,
        left: left * .5,
        name
      };
    });

    this.state = {
      mode: 'easy',
      population: this.getInitialPosition()
    };

    this.lock = false;

    this.mode = (mode) => {
      if (mode === 'current') return this.state.mode;
      this.setState({ mode });
      if (mode === 'easy') {
        this.setState({ population: this.getInitialPosition() });
      }
    };

    this.mouseMove = this.mouseMove.bind(this);
    this.runaway = this.runaway.bind(this);
    this.destination = this.destination.bind(this);
  }

  componentDidMount () {
    this.rect = this.container.getBoundingClientRect();
  }

  destination (width, height) {
    const minLeft = 10;
    const maxLeft = this.rect.width - 10 - width;
    const minTop = 10;
    const maxTop = this.rect.height - 10 - height;
    return {
      left: getRandomInt(minLeft, maxLeft),
      top: getRandomInt(minTop, maxTop)
    };
  }

  runaway (width, height, index) {
    return (event) => {
      const population = this.state.population.slice();
      const coords = this.destination(width, height);
      population[index].top = coords.top;
      population[index].left = coords.left;
      this.setState({ population });
    };
  }

  mouseMove (proxy, event) {
    if (this.lock) return;
    this.lock = true;
    this.forceUpdate();
    const { clientX, clientY } = proxy;
    this.setState(
      { clientX, clientY },
      () => {
        window.setTimeout(
          () => {
            this.lock = false;
            this.forceUpdate();
          },
          100
        );
      }
    );
  }

  render () {
    const {
      props,
      state,
      lock,
      mode,
      runaway
    } = this;

    const current = mode('current');
    const wait = lock || current === 'easy';

    return (
      <div
        className="homepage"
        ref={(el) => { this.container = el }}
        onMouseMove={wait ? null : this.mouseMove}
      >

        <Header
          mode={mode}
          isAdmin={props.isAdmin}
          path='/'
        />

        <Tapestry
          mode={mode}
          population={state.population}
          clientX={state.clientX}
          clientY={state.clientY}
          runaway={runaway}
        />

      </div>
    );
  }
}

function Header({ mode, isAdmin, path }) {
  return (
    <header>
      <Link text='Home' url='/' current={path} />
      <Link text='Blog' url='/blog' current={path} />
      {isAdmin && (
        <span>
          <Link text='Logout' url='/login' current={path} />
          <Link text='New blog post' url='/blog/post' current={path} />
          <Link text='Manage media' url='/media' current={path} />
        </span>
      )}
      <Toggler name='easy' mode={mode} />
      <Toggler name='hard' mode={mode} />
    </header>
  );
}

function Link ({ text, url, current }) {
  const disabled = url === current;
  return (
    <a
      className={`link d${disabled}`}
      href={url}
      onClick={disabled ? e => { e.preventDefault() } : null}
    >{text}</a>
  );
}

function Toggler ({ name, mode }) {
  const current = mode('current');
  const disabled = name === current;
  return (
    <a
      className={`${name} d${disabled}`}
      href={`#${name}`}
      onClick={(event) => {
        event.preventDefault();
        if (!disabled) mode(name);
      }}
    >{name} mode</a>
  );
}

function Tapestry ({ mode, population, clientX, clientY, runaway }) {
  return (
    <div style={{ position: 'relative' }}>
      {population.map((char, i) => (
        <Character
          clientX={clientX}
          clientY={clientY}
          runaway={runaway}
          key={i}
          index={i}
          {...char}
        />
      ))}
    </div>
  );
}

const px = num => `${num}px`;

function Character ({ name, width, height, left, top, clientX, clientY, runaway, index }) {
  if (clientX && clientY) {

  }
  return (
    <figure
      style={{
        position: "absolute",
        width: px(width),
        height: px(height),
        left: px(left),
        top: px(top),
        margin: '1em',
        transition: 'all .5s ease'
      }}
    >
      <img
        style={{
          position: "absolute",
          height: "100%",
          top: 0, left: 0, right: 0, bottom: 0
        }}
        src={`/static/images/${name}.png`}
        alt={name}
        onMouseEnter={runaway(width, height, index)}
      />
    </figure>
  );
}
