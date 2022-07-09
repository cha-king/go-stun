# Notes
- We need something running asynchronously that's populating stun and app buffers
- It needs to populate some sort of a ring buffer / FIFO
    - Two options: channels or byte slice
    - Byte slice is nice because it works well with reader interface
    - Also sucks for this because we somehow need to track addresses as well
    - Also, implementing a FIFO on a byte slice sucks
    - Channels also let us avoid needing mutexes as well too
- Stick with virtual conn type that sources data from channels
- Needs to be stateful
    - Handle when packet has only been partially read for given addr
    - [Example](https://stackoverflow.com/a/40643767/17926959)


# References
- [pfilter](https://github.com/AudriusButkevicius/pfilter)
