import React from "react";
import classNames from "classnames";
import PropTypes from "prop-types";
import {
  Col,
  Badge,
  FormCheckbox,
} from "shards-react";



class Toolbar extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      checked: false,
      callBack: props.callBack,
    };

    // This binding is necessary to make `this` work in the callback
    this.click = this.click.bind(this);
  }
  
  click(e) {
    this.setState(state => ({
      checked: !state.checked
    }));
    console.log(this.state.checked)
  };

  render() {
    const classes = classNames(
      "text-center",
      "text-md-left",
      "mb-sm-0",
    );

    return (
        <FormCheckbox className={classes} toggle onClick={this.click} onChange={this.state.callBack}>
          {this.state.checked ? 'Cards' : 'Table view'}
        </FormCheckbox>
    );
  }
};

export default Toolbar;