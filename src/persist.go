package movo

import (
    "appengine"
    "appengine/datastore"
    "appengine/memcache"
    "bytes"
    "gob"
    "os"
    "strconv"
    "time"
    )

/********************************************************
 * DSActDetail
 ********************************************************/
type DSActDetail struct {
  ActId int
  Type string
  Name string
  Desc string `datastore:",noindex"`
  Location string
  Start datastore.Time
  End datastore.Time
  NCadets int `datastore:",noindex"`
  NStaff int `datastore:",noindex"`
}

func (a *ActDetail) toDS() (res *DSActDetail) {
  res = &DSActDetail{
          ActId: a.ActId,
          Type: a.Type,
          Name: a.Name,
          Desc: a.Desc,
          Location: a.Location,
          NCadets: a.NCadets,
          NStaff: a.NStaff,
        }
  if a.Start != nil {
    res.Start = datastore.SecondsToTime(a.Start.Seconds())
  }
  if a.End != nil {
    res.End = datastore.SecondsToTime(a.End.Seconds())
  }

  return
}

func (d *DSActDetail) fromDS() (res *ActDetail) {
  res = &ActDetail{
          ActId: d.ActId,
          Type: d.Type,
          Name: d.Name,
          Desc: d.Desc,
          Location: d.Location,
          Start: d.Start.Time(),
          End: d.End.Time(),
          NCadets: d.NCadets,
          NStaff: d.NStaff,
        }

  return
}

func ActExists(c appengine.Context, actId int) (bool, os.Error) {
  return false, nil
}

func RetrieveActDetails(c appengine.Context, actId int) (res *ActDetail, err os.Error) {
  var d DSActDetail
  if itm, err := memcache.Get(c, "actId__" + strconv.Itoa(actId));
      err != nil && err != memcache.ErrCacheMiss {
    return nil, err
  } else if err == nil {
    // Cache hit
    buf := bytes.NewBuffer(itm.Value)
    dec := gob.NewDecoder(buf)
    dec.Decode(&d)
  } else {
    // Cache miss
    key := datastore.NewKey(c, "DSActDetail", "", int64(actId), nil)
    if err := datastore.Get(c, key, &d); err != nil {
      return nil, err
    }
    // XXX Add to cache
  }
  return d.fromDS(), nil
}

func (a *ActDetail) Persist(c appengine.Context) os.Error {
  d := a.toDS()
c.Debugf("Done converting to DSActDetail")
  key := datastore.NewKey(c, "DSActDetail", "", int64(d.ActId), nil)
c.Debugf("Done Creating Key")
  _,err := datastore.Put(c, key, d)
c.Debugf("Done persisting to datastore")
  if err != nil {
    return err
  }

  buf := bytes.NewBufferString("")
  enc := gob.NewEncoder(buf)
  enc.Encode(d)

  itm := &memcache.Item{
          Key:    "actId__" + strconv.Itoa(a.ActId),
          Value:  buf.Bytes(),
       }

  err = memcache.Set(c, itm)
c.Debugf("Request cache to memcache")

  return err
}

/********************************************************
 * ActiveAct
 ********************************************************/
type ActiveAct []int

var (
    activeForward = DayToSeconds(90)
    activeBehind = DayToSeconds(14)
    )

func UpdateActiveAct(c appengine.Context) os.Error {
  err := memcache.Delete(c, "activeAct")
  if err != nil {
    return err
  }
  _,err = GetActiveAct(c)
  return err
}

func GetActiveAct(c appengine.Context) ([]*ActDetail, os.Error) {
  var aa ActiveAct
  if itm, err := memcache.Get(c, "activeAct");
      err != nil && err != memcache.ErrCacheMiss {
    return nil, err
  } else if err == nil {
    // Cache hit
    buf := bytes.NewBuffer(itm.Value)
    dec := gob.NewDecoder(buf)
    dec.Decode(&aa)
  } else {
    // Cache miss
    q := datastore.NewQuery("DSActDetail").
         Filter("End >", datastore.SecondsToTime(time.Seconds() - activeBehind))

    ds := make([]DSActDetail, 0)
    if _,err = q.GetAll(c, &ds); err != nil {
      return nil, err
    }
    aa = make(ActiveAct, 0)
    for _,d := range ds {
      if d.Start.Time().Seconds() <= time.Seconds() + activeForward {
        aa = append(aa, d.ActId)
        buf := bytes.NewBufferString("")
        enc := gob.NewEncoder(buf)
        enc.Encode(d)

        itm := &memcache.Item{
                Key:    "actId__" + strconv.Itoa(d.ActId),
                Value:  buf.Bytes(),
             }

        err = memcache.Set(c, itm)
        c.Debugf("Request cache to memcache")
      }
    }

    buf := bytes.NewBufferString("")
    enc := gob.NewEncoder(buf)
    enc.Encode(aa)

    itm = &memcache.Item{
      Key: "activeAct",
      Value: buf.Bytes(),
      Expiration: int32(DayToSeconds(1)),
    }
    err = memcache.Set(c, itm)
    if err != nil {
      return nil, err
    }
  }

  var res []*ActDetail
  for _,id := range aa {
    a,err := RetrieveActDetails(c, id)
    if err != nil {
      return nil, err
    }
    res = append(res, a)
  }
  return res, nil
}

func init() {
  gob.Register(ActiveAct{})
  gob.Register(ActDetail{})
}
