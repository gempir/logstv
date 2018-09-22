import React, { Component } from 'react';
import moment from "moment";

export default class Filter extends Component {

    constructor(props) {
        super(props);

        this.state = {
            channel: "",
            username: "",
            year: moment().year(),
            month: moment().format("MMMM")
         }
    }

    render() {
		return (
            <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                <input
                    type="text"
                    className="channel-filter"
                    placeholder="channel"
                    onChange={this.onChannelChange}
                />
                <input
                    type="text"
                    className="username-filter"
                    placeholder="channel"
                    onChange={this.onUsernameChange}
                />
                <button type="submit" className="show-logs">show logs</button>
            </form>
		)
    }

    onChannelChange = (event) => {
        this.setState({...this.state, channel: event.target.value});
    } 

    onUsernameChange = (event) => {
        this.setState({...this.state, username: event.target.value});
    }
    
    onSubmit = (e) => {
        e.preventDefault();
        this.props.searchLogs(this.state.channel, this.state.username, this.state.year, this.state.month);
    }
}