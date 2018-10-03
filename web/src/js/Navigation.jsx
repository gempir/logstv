import React, { Component } from 'react';
import moment from "moment";

export default class Navigation extends Component {

    render() {
		return (
            <nav>
                <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                    <img src="/images/logo.png"></img>
                    <input
                        type="text"
                        className="channel-filter"
                        placeholder="channel"
                        spellcheck="false"
                    />
                    <input
                        type="text"
                        className="username-filter"
                        placeholder="username"
                        spellcheck="false"
                    />
                    <button type="submit" className="show-logs"><i class="material-icons">search</i></button>
                </form>
            </nav>
		)
    }
}