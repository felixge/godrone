/** @jsx React.DOM */
"use strict";

// Table implements a React HTML table that supports horizontal and vertical
// layout, as well as formatting of floating point numbers.
window.Table = (function() {
  return React.createClass({
    render: function() {
      return (
        <div>
          <h2>{this.props.title}</h2>
          {this.props.children}
          <table className={this.props.layout}>
            {
              this.props.layout == 'horizontal'
                ? this.renderHorizontal(this.props.data)
                : this.renderVertical(this.props.data)
            }
          </table>
        </div>
      );
    },
    renderHorizontal: function(data) {
      var headers = [];
      for (var key in data) {
        headers.push(<th>{key}</th>);
      }

      var values = [];
      for (var key in data) {
        var val = this.renderVal(data[key]);
        values.push(<td>{val}</td>);
      }
      return [
        <thead>
          <tr>{headers}</tr>
        </thead>,
        <tbody>
          <tr>{values}</tr>
        </tbody>,
      ];
    },
    renderVertical: function(data) {
      var rows = [];
      for (var key in data) {
        rows.push(
          <tr>
            <th>{key}</th>
            <td>{this.renderVal(data[key])}</td>
          </tr>
        );
      }

      return rows;
    },
    renderVal: function(val) {
      if (typeof val == 'number' && this.props.decimals) {
        val = (val||0).toFixed(this.props.decimals);
      }
      return val;
    },
  });
})();
