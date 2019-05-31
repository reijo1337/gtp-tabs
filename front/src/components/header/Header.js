import React, {Component} from 'react';
import TopHeader from "./TopHeader";
import Categories from "./Categories";
import SearchString from "./SearchString";

class Header extends Component {
    render() {
        return (
            <div>
                <TopHeader/>
                <Categories/>
                <SearchString/>
            </div>
        );
    }
}

export default Header;