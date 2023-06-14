'use client'
import { useEffect, useState } from 'react';


export default function Home() {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true)
    fetch(`http://localhost:8080/`, {mode: 'cors'})
    .then(response => response.json())
    .then(json => {
      //Handle filtering data from server
      const currentDate = new Date(Date.now());
      let pastDate = new Date();
      const inclusiveDate = pastDate.setDate(currentDate.getDate() - 2);
      const startDate = new Date(inclusiveDate);
      const filteredEntries = json.feed.entry.filter((entry) => new Date(Date.parse(entry.updated.label)) <= currentDate && new Date(Date.parse(entry.updated.label)) >= startDate);
      setReviews(filteredEntries)
    })
    .finally(setLoading(false))
  }, [])


  return (
    <div className='main'>
      {loading ? (<div>Loading...</div>) : (
        <>
        <h1>Reviews</h1>
        <table border={1}>
          <thead>
          <tr>
            <th>Username</th>
            <th>Review Content</th>
            <th>Rating</th>
            <th>Updated</th>
          </tr>
          </thead>
          {reviews.map(review => (
            <tbody key={review.author.uri.label}>
            <tr >
              <td>{review.author.name.label}</td>
              <td>{review.content.label}</td>
              <td>{review['im:rating'].label}</td>
              <td>{review.updated.label}</td>
            </tr>
            </tbody>
          ))}
        </table>
        </>
      )}
    </div>
  )
}