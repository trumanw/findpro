**Keywords**: indoor GPS, WiFi positioning, indoor mapping, indoor navigation, indoor positioning

# About

**The Framework for Internal Navigation and Discovery** (_FIND_) allows you to use your (Android) smartphone or WiFi-enabled computer (laptop or Raspberry Pi or etc.) to determine your position within your home or office. You can easily use this system in place of motion sensors as its resolution will allow your phone to distinguish whether you are in the living room, the kitchen or the bedroom, etc. The position information can then be used in a variety of ways including home automation, way-finding, or tracking! For more details, refer to [schollz/find](https://github.com/schollz/find).

Based on the concept of [schollz/find](https://github.com/schollz/find), [findpro](https://github.com/trumanw/findpro) re-factored the [find](https://github.com/schollz/find) framework to be feasible, reusable and production-ready from these aspects:

**Features**
- Re-write the positioning associated APIs with gRPC framework.
- Easier to launch a service via [cobra](https://github.com/spf13/cobra) commands.
- Add instruments interface which can be integrated with services like Prometheus, NewRelic, etc.
- Add new database interface for more KV database like [etcd](https://github.com/coreos/etcd).
- Add new cache interface for [Redis](https://github.com/go-redis/redis).
- Docker compose support.
- Support classifier written in Python. ([google/grumpy](https://github.com/google/grumpy))
- Add integration with 3rd-party libraries like [scikit](http://scikit-learn.org/), TensorFlow, etc.
