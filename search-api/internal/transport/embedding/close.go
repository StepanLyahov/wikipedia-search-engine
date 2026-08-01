package embedding

// Close releases the underlying gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
