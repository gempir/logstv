import React, { Component } from 'react';
import moment from 'moment';

export default class LogView extends Component {
    render() {

		return (
			<div className="log-view">
				{/* {this.props.logs.map((value, key) => 
					<div key={key} className="line">
						<a href={`#${value.timestamp}`}><span className="timestamp">{this.formatDate(value.timestamp)}</span></a> {value.text}
					</div>
				)} */}
			</div>
		);
	}
	
	formatDate = (timestamp) => {
		return moment(timestamp).format("YYYY-MM-DD HH:mm:ss UTC");
	}
}