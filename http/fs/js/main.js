/** @jsx React.DOM */
"use strict";

(function() {
  var GoDrone = React.createClass({
    getInitialState: function() {
      return {
        wsStatus: 'Disconnected',
        wsMsgsPerSec: 0,
        gamepad: {Status: 'Disconnected'},
        sensors: {},
        attitude: {},
        raw: {},
        msgsPerSec: 0,
      };
    },
    componentDidMount: function() {
      this._handleGamepad();
      this._handleClient();
    },
    _handleClient: function() {
      var self = this;
      var msgsPerSec = 0;
      setInterval(function() {
        self.setState({wsMsgsPerSec: msgsPerSec});
        msgsPerSec = 0;
      }, 1000);

      var firstConnect = true;
      (new Client({
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
      })).connect();
    },
    _handleGamepad: function() {
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
        onChange: function(rawState) {
          var state = {};
          rawState.axes.forEach(function(val, i) {
            state['A'+i] = val;
          });
          rawState.buttons.forEach(function(val, i) {
            state['B'+i] = val;
          });

          var gamepad = $.extend(true, {}, self.state.gamepad, state);
          self.setState({gamepad: gamepad});
        },
      })).connect();
    },
    render: function() {
      return (
        <div>
          <h1>GoDrone</h1>

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
              right now, but should work with others. B0 (X) square toggles the
              motors on/off. Right joystick controls throttle, left joystick
              controls movement. There is no rotation yet. You may also
              experience buttons that keep toggling as you hold them down, in
              this case reconnect your controller. It seems to be a Gamepad API
              issue.
            </p>
          </Table>
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

  React.renderComponent(<GoDrone version="0.1" wsUrl="ws://192.168.1.1/ws"/>, document.body);
})();
