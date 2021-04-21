import React from 'react';
import PageTypes from '../../../../Constants/PageTypes/PageTypes';

type CardProps = {
  title: string,
  paragraph: string
}

const Card = ({ cardInfo, setPage }: CardProps) => 
  <div className="card-div" onClick={(e) => {setPage(e, PageTypes.clean)
                                              localStorage.setItem('dash', JSON.stringify(cardInfo))}}>
    <div className="card">
        <img className="pic" src="img/aesthetic-anime-13.jpg" alt="abstract art"></img>
        <div className="cardDesc">
          <h1>{cardInfo.creator.userName}'s Dash</h1>
        </div>
    </div>
  </div>

export default Card;