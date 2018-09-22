import React, { Component } from "react";
import "whatwg-fetch";
import Filter from "./Filter";
import LogView from "./LogView";
import moment from "moment";

export default class LogSearch extends Component {
	constructor(props) {
        super(props);	
        
        this.state = {
            channels: [],
            logs: [],
            visibleLogs: [],
            isLoading: false,
        };

		fetch("https://api.gempir.com/channel").then((response) => {
			return response.json()
		}).then((json) => {
			this.setState({...this.state, channels: json.channels, selectedChannel: json.channels[0]});
        });
	}

	render() {
		return (
			<div className="log-search">
                <Filter 
                    channels={this.state.channels} 
                    searchLogs={this.searchLogs} 
                /> 
                <LogView logs={this.state.visibleLogs} isLoading={this.state.isLoading}/>
			</div>
		);
    }

    searchLogs = (channel, username, year, month) => {
        this.setState({...this.state, isLoading: true});

        let options = {
            headers: {
                "Content-Type": "application/json"
            }
        }

        fetch(`https://api.gempir.com/channel/${channel}/user/${username}/${year}/${month}`, options).then(this.checkStatus).then((response) => {
			return response.json()
		}).then((json) => {

            for (let value of json) {
                value.timestamp = Date.parse(value.timestamp)
            }

            this.setState({...this.state, isLoading: false, logs: json, visibleLogs: json});
		}).catch((error) => {
            this.setState({...this.state, isLoading: false, logs: [], visibleLogs: []});
        });
    }

    checkStatus = (response) => {
        if (response.status >= 200 && response.status < 300) {
            return response
        } else {
            var error = new Error(response.statusText)
            error.response = response
            throw error
        }
    }
}