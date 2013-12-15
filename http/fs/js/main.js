/** @jsx React.DOM */
"use strict";

(function() {
  var GoDrone = React.createClass({
    getInitialState: function() {
      return {
        wsStatus: 'Disconnected',
        wsMsgsPerSec: 0,
        gamepad: {Status: 'Disconnected'},
        fly: false,
        control: {Roll: 0, Pitch: 0, Yaw: 0, Throttle: 0},
        sensors: {},
        attitude: {},
        raw: {},
        msgsPerSec: 0,
      };
    },
    componentDidMount: function() {
      var client = this._handleClient();
      this._handleGamepad(client);
    },
    _handleClient: function() {
      var self = this;
      var msgsPerSec = 0;
      setInterval(function() {
        self.setState({wsMsgsPerSec: msgsPerSec});
        msgsPerSec = 0;
      }, 1000);

      var firstConnect = true;
      var client = new Client({
        url: this.props.wsUrl,
        onConnecting: function() {
          self.setState({wsStatus: 'Connecting'});
        },
        onConnect: function() {
          if (firstConnect) {
            self.setState({wsStatus: 'Connected'});
          } else {
            self.setState({wsStatus: 'Reloading page'});
            setTimeout(function() {
              window.location.reload(true);
            }, 500);
          }
          firstConnect = false;
        },
        onError: function(err) {
          self.setState({wsStatus: 'Error'});
        },
        onClose: function(e) {
          self.setState({wsStatus: 'Disconnected'});
        },
        onData: function(data) {
          msgsPerSec++;

          var sensors = {};
          ['Ax', 'Ay', 'Az', 'Gx', 'Gy', 'Gz'].forEach(function(key) {
            sensors[key] = data.NavData[key];
          });
          var attitude = {};
          ['Roll', 'Pitch', 'Yaw', 'Altitude'].forEach(function(key) {
            attitude[key] = data.AttitudeData[key];
          });
          self.setState({
            sensors: sensors,
            attitude: attitude,
            raw: data.NavData.Raw,
          });
        },
      });
      client.connect();
      return client;
    },
    _handleGamepad: function(client) {
      var fly = false;
      var prevGamepad;

      var self = this;
      (new Gamepad({
        onConnect: function() {
          var gamepad = $.extend(true, {}, self.state.gamepad);
          gamepad.Status = 'Connected';
          self.setState({gamepad: gamepad});
        },
        onClose: function() {
          var gamepad = $.extend(true, {}, self.state.gamepad);
          gamepad.Status = 'Disconnected';
          self.setState({gamepad: gamepad});
        },
        onChange: function(gamepad) {
          var controlState = $.extend(true, {}, self.state.control);
          var gamepadState = $.extend(true, {}, self.state.gamepad);
          var throttle = 0;
          var roll = 0;
          var pitch = 0;

          function axisToAngle(val) {
            if (Math.abs(val) < self.props.gamepadAxisMin) {
              return 0;
            }
            return val * self.props.maxAngle;
          }

          gamepad.axes.forEach(function(val, i) {
            gamepadState['A'+i] = val;

            if (i === self.props.throttleAxis) {
              throttle = Math.max(0, -val);
            } else if (i === self.props.pitchAxis) {
              pitch = axisToAngle(val);
            } else if (i === self.props.rollAxis) {
              roll = axisToAngle(val);
            } 
          });
          gamepad.buttons.forEach(function(val, i) {
            gamepadState['B'+i] = val;

            var prevVal = (prevGamepad && prevGamepad.buttons[i]);
            if (i === self.props.flyButton && val === 0 && prevVal === 1) {
              fly = !fly;
            }
          });

          if (!fly) {
            controlState.Throttle = 0;
          } else {
            controlState.Throttle = (self.props.maxThrottle-self.props.flyThrottle)*throttle + self.props.flyThrottle;
          }
          controlState.Pitch = pitch;
          controlState.Roll = roll;

          client.send(controlState);
          self.setState({
            gamepad: gamepadState,
            control: controlState,
          });
          prevGamepad = gamepad;
        },
        onButtonPress: function(button) {
          if (button !== self.props.flyButton) {
            return;
          }

          self._fly = !self._fly;
          self._updateControl();
        },
      })).connect();
    },
    _updateControl: function() {
      var control = $.extend(true, {}, this.state.control);
      if (!this._fly) {
        control.Throttle = 0;
      } else {
        control.Throttle = this.props.flyThrottle;
      }
      this.setState({control: control});
    },
    render: function() {
      return (
        <div>
          <h1>GoDrone {this.props.version}</h1>

          <Table
            title="WebSocket"
            layout="horizontal"
            data={{Status: this.state.wsStatus, MsgsPerSec: this.state.wsMsgsPerSec}}
          />
          <Table
            title="Gamepad"
            layout="horizontal"
            decimals={2}
            data={this.state.gamepad}
          >
            <p>
              Only tested with <a href="http://goo.gl/OaDB9J">this</a> gamepad
              right now, but should work with others. B{this.props.flyButton}
              toggles the motors on/off. Right joystick controls throttle, left
              joystick controls movement. There is no rotation yet. You may
              also experience buttons that keep toggling as you hold them down,
              in this case reconnect your controller. It seems to be a Gamepad
              API issue.
            </p>
          </Table>
          <Table
            title="Control"
            layout="horizontal"
            decimals={2}
            data={this.state.control}
          />
          <Table
            title="Attitude"
            layout="horizontal"
            decimals={2}
            data={this.state.attitude}
          >
            <p>
              Describes the drone's position relative to the ground. See
              <a href="http://en.wikipedia.org/wiki/Aircraft_principal_axes">
              Aircraft principal axes.</a> for more information.
            </p>
          </Table>
          <Table
            title="Sensor Data"
            layout="horizontal"
            decimals={2}
            data={this.state.sensors}
          >
            <p>
              The sensor data, adjusted for
              <a href="http://en.wikipedia.org/wiki/Sensor#Sensor_deviations">
              bias and sensitivity</a>. Units are m/s² (accelerometers) °/s
              (gyroscopes).
            </p>
          </Table>
          <Table
            layout="vertical"
            title="RAW Navboard Data"
            data={this.state.raw}
          >
            <p>The raw data decoded from the navboard tty file.</p>
          </Table>
        </div>
      );
    },
  });

  React.renderComponent(
    <GoDrone
      version="0.1"
      wsUrl="ws://192.168.1.1/ws"
      gamepadAxisMin={0.01}
      maxAngle={3}
      flyButton={0}
      pitchAxis={1}
      rollAxis={0}
      throttleAxis={3}
      flyThrottle={0.4}
      maxThrottle={1}
    />,
    document.body
  );
})();
