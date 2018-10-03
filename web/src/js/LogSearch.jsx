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
	}

	render() {
		return (
			<div className="log-search">
                <Filter 
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

        fetch(`https://api.logs.tv/channel/${channel}/user/${username}`, options).then(this.checkStatus).then((response) => {
			return response.json()
		}).then((json) => {
        
            console.log(json);

            for (let value of json.messages) {
                value.timestamp = Date.parse(value.timestamp)
            }

            this.setState({...this.state, isLoading: false, logs: json.messages, visibleLogs: json.messages});
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